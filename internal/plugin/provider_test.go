// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2026 The semrel Authors

package plugin

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

type commandCall struct {
	Name string
	Args []string
	Dir  string
}

type fakeRunner struct {
	Calls []commandCall
}

func (f *fakeRunner) Run(_ context.Context, name string, args []string, _ []string, dir string) error {
	copiedArgs := append([]string(nil), args...)
	f.Calls = append(f.Calls, commandCall{Name: name, Args: copiedArgs, Dir: dir})
	return nil
}

func TestMavenUpdaterExecuteRunsExpectedCommands(t *testing.T) {
	t.Parallel()

	runner := &fakeRunner{}
	updater := NewMavenUpdater(".")
	updater.Settings = "settings.xml"
	updater.Runner = runner

	result, err := updater.Execute(context.Background(), &Release{Version: "2.0.0"})
	require.NoError(t, err)
	require.Len(t, runner.Calls, 2)
	require.Equal(t, []string{"-B", "--settings", "settings.xml", "versions:set", "-DnewVersion=2.0.0", "-DgenerateBackupPoms=false"}, runner.Calls[0].Args)
	require.Equal(t, []string{"-B", "--settings", "settings.xml", "deploy", "-DskipTests"}, runner.Calls[1].Args)
	require.Equal(t, "2.0.0", result.Outputs["version"])
}
