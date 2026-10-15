package ui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/ansi"
	"github.com/sudabon/medit/internal/loader"
)

func TestFindMatchesSmartCase(t *testing.T) {
	lines := []string{"Foo foo", "bar", "\x1b[1mFOO\x1b[0m"}

	got := findMatches(lines, "foo")
	want := []searchMatch{{0, 0, 3}, {0, 4, 7}, {2, 0, 3}}
	if len(got) != len(want) {
		t.Fatalf("matches = %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("matches[%d] = %v, want %v", i, got[i], want[i])
		}
	}

	if got := findMatches(lines, "Foo"); len(got) != 1 || got[0] != (searchMatch{0, 0, 3}) {
		t.Fatalf("case-sensitive matches = %v", got)
	}
}

func TestFindMatchesWideChars(t *testing.T) {
	got := findMatches([]string{"  日本語の検索テスト"}, "検索")
	if len(got) != 1 || got[0] != (searchMatch{0, 10, 14}) {
		t.Fatalf("matches = %v, want [{0 10 14}]", got)
	}
}

func TestHighlightLinePreservesText(t *testing.T) {
	line := "\x1b[31mhello\x1b[0m 日本語 world"
	ms := findMatches([]string{line}, "日本")
	out := highlightLine(line, ms, 0)
	if ansi.Strip(out) != ansi.Strip(line) {
		t.Fatalf("visible text changed: %q -> %q", ansi.Strip(line), ansi.Strip(out))
	}
}

func TestSearchKeysJumpAndCycle(t *testing.T) {
	content := "# Title\n\nalpha\n\n"
	for range 60 {
		content += "filler paragraph\n\n"
	}
	content += "needle one\n\n"
	for range 60 {
		content += "filler paragraph\n\n"
	}
	content += "needle two\n"

	var tm tea.Model = New(&loader.Document{Path: "x.md", Content: content})
	tm, _ = tm.Update(tea.WindowSizeMsg{Width: 100, Height: 20})

	press := func(keys ...string) {
		for _, k := range keys {
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(k)}
			switch k {
			case "enter":
				msg = tea.KeyMsg{Type: tea.KeyEnter}
			case "esc":
				msg = tea.KeyMsg{Type: tea.KeyEsc}
			}
			tm, _ = tm.Update(msg)
		}
	}

	press("s", "n", "e", "e", "d", "l", "e", "enter")
	m := tm.(Model)
	if m.searching {
		t.Fatal("still in search mode after enter")
	}
	if len(m.matches) != 2 || m.matchIdx != 0 {
		t.Fatalf("matches = %v, idx = %d", m.matches, m.matchIdx)
	}
	first := m.matches[0].line
	if first < m.viewport.YOffset || first >= m.viewport.YOffset+m.viewport.Height {
		t.Fatalf("match line %d not visible (offset %d)", first, m.viewport.YOffset)
	}

	press("n")
	if m = tm.(Model); m.matchIdx != 1 {
		t.Fatalf("after n idx = %d, want 1", m.matchIdx)
	}
	press("n")
	if m = tm.(Model); m.matchIdx != 0 {
		t.Fatalf("after wrap idx = %d, want 0", m.matchIdx)
	}
	press("p")
	if m = tm.(Model); m.matchIdx != 1 {
		t.Fatalf("after p idx = %d, want 1", m.matchIdx)
	}

	press("esc")
	if m = tm.(Model); m.query != "" || len(m.matches) != 0 {
		t.Fatalf("esc did not clear search: %q %v", m.query, m.matches)
	}
}

func TestSearchNotFoundAndCancel(t *testing.T) {
	var tm tea.Model = New(&loader.Document{Path: "x.md", Content: "hello\n"})
	tm, _ = tm.Update(tea.WindowSizeMsg{Width: 80, Height: 10})

	tm, _ = tm.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("s")})
	tm, _ = tm.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("zzz")})
	tm, _ = tm.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m := tm.(Model)
	if m.message != "not found: zzz" || m.query != "" {
		t.Fatalf("message = %q, query = %q", m.message, m.query)
	}

	tm, _ = tm.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("s")})
	tm, _ = tm.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("q")})
	if m = tm.(Model); !m.searching || m.input.Value() != "q" {
		t.Fatalf("q should be typed into the prompt, got searching=%v value=%q", m.searching, m.input.Value())
	}
	tm, _ = tm.Update(tea.KeyMsg{Type: tea.KeyEsc})
	if m = tm.(Model); m.searching || m.query != "" {
		t.Fatalf("esc did not cancel input")
	}
}
