// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2026 The semrel Authors

// Package plugin updates Maven pom.xml files in-place.
package plugin

import (
	"fmt"
	"os"
	"regexp"
)

var pomVersionPattern = regexp.MustCompile(`(?s)(<project[^>]*>.*?<version>)([^<]*)(</version>)`)

// Updater updates pom.xml version fields.
type Updater struct{}

// NewUpdater creates an updater.
func NewUpdater() *Updater {
	return &Updater{}
}

// Update rewrites the project version in pom.xml.
func (u *Updater) Update(path, version string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read %s: %w", path, err)
	}
	if !pomVersionPattern.Match(data) {
		return fmt.Errorf("project version not found in %s", path)
	}
	updated := pomVersionPattern.ReplaceAllString(string(data), `${1}`+version+`${3}`)
	if err := os.WriteFile(path, []byte(updated), 0o644); err != nil {
		return fmt.Errorf("write %s: %w", path, err)
	}
	return nil
}
