package loader

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadExistingFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "note.md")
	content := "# Hello\n日本語\n"
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	doc, err := Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if doc.IsNew {
		t.Fatal("expected existing file, got IsNew")
	}
	if doc.Content != content {
		t.Fatalf("content = %q, want %q", doc.Content, content)
	}
	if doc.Path != path {
		t.Fatalf("path = %q, want %q", doc.Path, path)
	}
}

func TestLoadNonexistentAsNew(t *testing.T) {
	path := filepath.Join(t.TempDir(), "missing.md")

	doc, err := Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if !doc.IsNew {
		t.Fatal("expected IsNew")
	}
	if doc.Content != "" {
		t.Fatalf("content = %q, want empty", doc.Content)
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatalf("file should not exist on disk, stat err = %v", err)
	}
}

func TestLoadDirectoryError(t *testing.T) {
	dir := t.TempDir()
	_, err := Load(dir)
	if err == nil {
		t.Fatal("expected error for directory")
	}
}

func TestLoadPermissionError(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "secret.md")
	if err := os.WriteFile(path, []byte("secret"), 0o000); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chmod(path, 0o644) })

	_, err := Load(path)
	if err == nil {
		t.Fatal("expected permission error")
	}
}

func TestReloadUpdatesContent(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "note.md")
	if err := os.WriteFile(path, []byte("v1"), 0o644); err != nil {
		t.Fatal(err)
	}
	doc, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(path, []byte("v2"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := doc.Reload(); err != nil {
		t.Fatalf("Reload: %v", err)
	}
	if doc.Content != "v2" {
		t.Fatalf("content = %q, want v2", doc.Content)
	}
}

func TestReloadMissingKeepsContent(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "note.md")
	if err := os.WriteFile(path, []byte("keep me"), 0o644); err != nil {
		t.Fatal(err)
	}
	doc, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(path); err != nil {
		t.Fatal(err)
	}

	err = doc.Reload()
	if err == nil {
		t.Fatal("expected reload error")
	}
	if doc.Content != "keep me" {
		t.Fatalf("content changed to %q", doc.Content)
	}
}

func TestReloadCreatesFromNew(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "new.md")
	doc, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if !doc.IsNew {
		t.Fatal("expected IsNew")
	}

	if err := os.WriteFile(path, []byte("# created"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := doc.Reload(); err != nil {
		t.Fatalf("Reload: %v", err)
	}
	if doc.IsNew {
		t.Fatal("expected IsNew cleared")
	}
	if doc.Content != "# created" {
		t.Fatalf("content = %q", doc.Content)
	}
}
