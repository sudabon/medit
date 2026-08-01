package render

import (
	"strings"
	"testing"
)

func TestRenderHeadingAndTable(t *testing.T) {
	r := New()
	r.SetTheme(ThemeDark)
	md := "# Title\n\n## Sub\n\n| A | B |\n|---|---|\n| 1 | 2 |\n"
	out, err := r.Render(md, 80)
	if err != nil {
		t.Fatalf("Render: %v", err)
	}
	if !strings.Contains(out, "Title") {
		t.Fatalf("missing heading text in output: %q", out)
	}
	if !strings.Contains(out, "Sub") {
		t.Fatalf("missing h2 text in output: %q", out)
	}
	plain := stripANSI(out)
	if strings.Contains(plain, "# Title") || strings.Contains(plain, "## Sub") {
		t.Fatalf("raw markdown heading markers should not appear: %q", plain)
	}
	if !strings.Contains(out, "A") || !strings.Contains(out, "B") {
		t.Fatalf("missing table cells in output: %q", out)
	}
}

func TestRenderLightThemeHasANSI(t *testing.T) {
	r := New()
	r.SetTheme(ThemeLight)
	out, err := r.Render("# Hello\n\nbody text\n", 60)
	if err != nil {
		t.Fatal(err)
	}
	plain := stripANSI(out)
	if strings.Contains(plain, "# Hello") {
		t.Fatalf("raw hash in light theme: %q", plain)
	}
	if !strings.Contains(out, "\x1b[") {
		t.Fatal("expected ANSI codes in light theme")
	}
}

func stripANSI(s string) string {
	var b strings.Builder
	for i := 0; i < len(s); i++ {
		if s[i] == 0x1b {
			for i < len(s) && s[i] != 'm' {
				i++
			}
			continue
		}
		b.WriteByte(s[i])
	}
	return b.String()
}

func TestRenderJapaneseWrap(t *testing.T) {
	md := "## 見出し\n\nこれは日本語の段落です。端末幅に合わせて折り返されることを確認するための長い文章を続けます。" +
		"全角文字の幅が正しく扱われ、表示が崩れないことが重要です。\n"
	r := New()
	out, err := r.Render(md, 40)
	if err != nil {
		t.Fatalf("Render: %v", err)
	}
	if !strings.Contains(out, "見出し") {
		t.Fatalf("missing Japanese heading: %q", out)
	}
	if !strings.Contains(out, "日本語") {
		t.Fatalf("missing Japanese body: %q", out)
	}
	// Narrow wrap should produce multiple lines.
	if strings.Count(out, "\n") < 2 {
		t.Fatalf("expected wrapped multiline output, got %q", out)
	}
}

func TestRenderCachesUntilChange(t *testing.T) {
	r := New()
	md := "# A\n"
	out1, err := r.Render(md, 60)
	if err != nil {
		t.Fatal(err)
	}
	out2, err := r.Render(md, 60)
	if err != nil {
		t.Fatal(err)
	}
	if out1 != out2 {
		t.Fatal("cached render should be identical")
	}

	if _, err := r.Render(md, 40); err != nil {
		t.Fatal(err)
	}
	out3, err := r.Render("# B\n", 60)
	if err != nil {
		t.Fatal(err)
	}
	if out3 == out1 {
		t.Fatal("content change should produce different output")
	}
}
