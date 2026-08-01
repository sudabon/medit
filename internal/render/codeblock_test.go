package render

import (
	"strings"
	"testing"
)

func TestSplitFenced(t *testing.T) {
	src := "# Hi\n\n```go\nfmt.Println(\"x\")\n```\n\nbye\n"
	segs := splitFenced(src)
	if len(segs) != 3 {
		t.Fatalf("segs = %d %#v", len(segs), segs)
	}
	if segs[0].isCode || !strings.Contains(segs[0].markdown, "# Hi") {
		t.Fatalf("seg0 = %#v", segs[0])
	}
	if !segs[1].isCode || segs[1].lang != "go" || !strings.Contains(segs[1].code, "Println") {
		t.Fatalf("seg1 = %#v", segs[1])
	}
	if segs[2].isCode || !strings.Contains(segs[2].markdown, "bye") {
		t.Fatalf("seg2 = %#v", segs[2])
	}
}

func TestRenderBoxedCodeBlock(t *testing.T) {
	r := New()
	r.SetTheme(ThemeLight)
	md := "before\n\n```bash\nexport VISUAL=nvim\n```\n\nafter\n"
	out, err := r.Render(md, 60)
	if err != nil {
		t.Fatal(err)
	}
	plain := stripANSI(out)
	if !strings.Contains(plain, "╭") || !strings.Contains(plain, "╰") {
		t.Fatalf("expected box borders in output: %q", plain)
	}
	if !strings.Contains(plain, "export") {
		t.Fatalf("missing code content: %q", plain)
	}
	if strings.Contains(plain, "```") {
		t.Fatalf("raw fences should not appear: %q", plain)
	}
}
