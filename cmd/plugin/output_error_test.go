package main

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type failingWriter struct{}

func (failingWriter) Write([]byte) (int, error) { return 0, errors.New("write failed") }

func TestRunSucceedsAfterMutationWhenOutputFails(t *testing.T) {
	file := filepath.Join(t.TempDir(), "pom.xml")
	if err := os.WriteFile(file, []byte("<project><version>1.0.0</version></project>"), 0o644); err != nil {
		t.Fatal(err)
	}
	env := map[string]string{"SEMREL_VERSION": "v1.1.0", "SEMREL_PLUGIN_FILE": file}
	if code := run(failingWriter{}, failingWriter{}, func(key string) string { return env[key] }); code != 0 {
		t.Fatalf("run() code = %d, want 0", code)
	}
	got, err := os.ReadFile(file)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(got), "<version>1.1.0</version>") {
		t.Fatalf("updated file = %s", got)
	}
}
