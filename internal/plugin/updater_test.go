package plugin

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestUpdaterUpdatePom(t *testing.T) {
	t.Parallel()

	dir, err := os.MkdirTemp("", "updater-maven-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	file := filepath.Join(dir, "pom.xml")
	original := "<project><modelVersion>4.0.0</modelVersion><groupId>x</groupId><artifactId>demo</artifactId><version>1.2.3</version></project>"
	if err := os.WriteFile(file, []byte(original), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := NewUpdater().Update(file, "1.3.0"); err != nil {
		t.Fatalf("Update() error = %v", err)
	}

	got, err := os.ReadFile(file)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(got), "<version>1.3.0</version>") {
		t.Fatalf("updated file = %s", got)
	}
}

func TestUpdaterMissingFile(t *testing.T) {
	t.Parallel()

	err := NewUpdater().Update(filepath.Join(t.TempDir(), "pom.xml"), "1.3.0")
	if err == nil || !strings.Contains(err.Error(), "read") {
		t.Fatalf("expected read error, got %v", err)
	}
}

func TestUpdaterMissingVersion(t *testing.T) {
	t.Parallel()

	dir, err := os.MkdirTemp("", "updater-maven-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	file := filepath.Join(dir, "pom.xml")
	if err := os.WriteFile(file, []byte("<project></project>"), 0o644); err != nil {
		t.Fatal(err)
	}

	err = NewUpdater().Update(file, "1.3.0")
	if err == nil || !strings.Contains(err.Error(), "project version not found") {
		t.Fatalf("expected version error, got %v", err)
	}
}
