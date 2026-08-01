package main

import (
	"strings"
	"testing"
)

func TestRunUsage(t *testing.T) {
	err := run(nil)
	if err == nil || !strings.Contains(err.Error(), "usage:") {
		t.Fatalf("got %v, want usage error", err)
	}
	err = run([]string{"a", "b"})
	if err == nil || !strings.Contains(err.Error(), "usage:") {
		t.Fatalf("got %v, want usage error", err)
	}
}

func TestRunDirectoryError(t *testing.T) {
	err := run([]string{t.TempDir()})
	if err == nil {
		t.Fatal("expected directory error")
	}
}
