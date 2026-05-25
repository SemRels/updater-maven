// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2026 The semrel Authors

package plugin

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Release contains the SemRel release data consumed by this plugin.
type Release struct {
	Version         string
	PreviousVersion string
	TagName         string
	Repository      string
	Changelog       string
	CommitSHA       string
	DryRun          bool
	Metadata        map[string]string
	Commits         []string
}

// Result captures the outcome of a plugin execution.
type Result struct {
	Name       string
	Outputs    map[string]string
	Skipped    bool
	SkipReason string
}

// Provider is the contract exposed by this plugin implementation.
type Provider interface {
	Name() string
	HealthCheck(context.Context) error
	Validate(map[string]interface{}) error
	Execute(context.Context, *Release) (*Result, error)
	ReleaseContext() []string
}

// CommandRunner executes external commands.
type CommandRunner interface {
	Run(context.Context, string, []string, []string, string) error
}

// ExecRunner runs external commands with os/exec.
type ExecRunner struct{}

func (ExecRunner) Run(ctx context.Context, name string, args []string, env []string, dir string) error {
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), env...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// MavenUpdater bumps the pom version and deploys artifacts.
type MavenUpdater struct {
	WorkingDir string
	Settings   string
	Runner     CommandRunner
}

// NewMavenUpdater constructs a Maven updater.
func NewMavenUpdater(workingDir string) *MavenUpdater {
	if strings.TrimSpace(workingDir) == "" {
		workingDir = "."
	}
	return &MavenUpdater{WorkingDir: workingDir, Settings: strings.TrimSpace(os.Getenv("MAVEN_SETTINGS")), Runner: ExecRunner{}}
}

func (m *MavenUpdater) Name() string { return "updater-maven" }

func (m *MavenUpdater) HealthCheck(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return nil
	}
}

func (m *MavenUpdater) Validate(map[string]interface{}) error {
	if strings.TrimSpace(m.WorkingDir) == "" {
		return fmt.Errorf("maven: working directory must not be empty")
	}
	if m.Runner == nil {
		m.Runner = ExecRunner{}
	}
	return nil
}

func (m *MavenUpdater) ReleaseContext() []string {
	return []string{"version"}
}

func (m *MavenUpdater) Execute(ctx context.Context, rel *Release) (*Result, error) {
	if err := m.HealthCheck(ctx); err != nil {
		return nil, err
	}
	if err := m.Validate(nil); err != nil {
		return nil, err
	}
	if rel == nil {
		return nil, fmt.Errorf("maven: release is required")
	}
	if strings.TrimSpace(rel.Version) == "" {
		return nil, fmt.Errorf("maven: release version is required")
	}
	dir, err := filepath.Abs(m.WorkingDir)
	if err != nil {
		return nil, fmt.Errorf("maven: resolve working dir: %w", err)
	}
	if rel.DryRun {
		return &Result{Name: m.Name(), Outputs: map[string]string{"working_dir": dir, "version": rel.Version, "dry_run": "true"}}, nil
	}

	baseArgs := []string{"-B"}
	if m.Settings != "" {
		baseArgs = append(baseArgs, "--settings", m.Settings)
	}

	setVersionArgs := append(append([]string{}, baseArgs...), "versions:set", "-DnewVersion="+rel.Version, "-DgenerateBackupPoms=false")
	if err := m.Runner.Run(ctx, "mvn", setVersionArgs, nil, dir); err != nil {
		return nil, fmt.Errorf("maven: versions:set failed: %w", err)
	}

	deployArgs := append(append([]string{}, baseArgs...), "deploy", "-DskipTests")
	if err := m.Runner.Run(ctx, "mvn", deployArgs, nil, dir); err != nil {
		return nil, fmt.Errorf("maven: deploy failed: %w", err)
	}

	return &Result{Name: m.Name(), Outputs: map[string]string{"working_dir": dir, "version": rel.Version}}, nil
}
