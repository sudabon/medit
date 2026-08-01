package ui

import "testing"

func TestResolveEditorDefault(t *testing.T) {
	t.Setenv("VISUAL", "")
	t.Setenv("EDITOR", "")
	cmd, args := resolveEditor()
	if cmd != "vi" {
		t.Fatalf("cmd = %q, want vi", cmd)
	}
	if len(args) != 0 {
		t.Fatalf("args = %v, want empty", args)
	}
}

func TestResolveEditorVisual(t *testing.T) {
	t.Setenv("VISUAL", "nvim -p")
	t.Setenv("EDITOR", "nano")
	cmd, args := resolveEditor()
	if cmd != "nvim" {
		t.Fatalf("cmd = %q, want nvim", cmd)
	}
	if len(args) != 1 || args[0] != "-p" {
		t.Fatalf("args = %v, want [-p]", args)
	}
}

func TestResolveEditorFallback(t *testing.T) {
	t.Setenv("VISUAL", "")
	t.Setenv("EDITOR", "vim")
	cmd, args := resolveEditor()
	if cmd != "vim" {
		t.Fatalf("cmd = %q, want vim", cmd)
	}
	if len(args) != 0 {
		t.Fatalf("args = %v", args)
	}
}
