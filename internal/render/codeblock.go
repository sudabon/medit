package render

import (
	"bytes"
	"strings"
	"unicode"

	"github.com/alecthomas/chroma/v2/quick"
	"github.com/charmbracelet/lipgloss"
)

type segment struct {
	markdown string
	code     string
	lang     string
	isCode   bool
}

// splitFenced splits Markdown into prose segments and fenced code blocks.
func splitFenced(src string) []segment {
	lines := strings.Split(src, "\n")
	var segs []segment
	var prose strings.Builder

	flushProse := func() {
		if prose.Len() == 0 {
			return
		}
		segs = append(segs, segment{markdown: prose.String()})
		prose.Reset()
	}

	for i := 0; i < len(lines); {
		line := lines[i]
		marker, lang, ok := parseFenceOpen(line)
		if !ok {
			if prose.Len() > 0 {
				prose.WriteByte('\n')
			}
			prose.WriteString(line)
			i++
			continue
		}

		flushProse()
		i++
		var code strings.Builder
		closed := false
		for i < len(lines) {
			if isFenceClose(lines[i], marker) {
				closed = true
				i++
				break
			}
			if code.Len() > 0 {
				code.WriteByte('\n')
			}
			code.WriteString(lines[i])
			i++
		}
		if !closed {
			// Unclosed fence: treat as prose including the opener.
			var raw strings.Builder
			raw.WriteString(marker)
			if lang != "" {
				raw.WriteString(lang)
			}
			raw.WriteByte('\n')
			raw.WriteString(code.String())
			segs = append(segs, segment{markdown: raw.String()})
			continue
		}
		segs = append(segs, segment{isCode: true, code: code.String(), lang: lang})
	}
	flushProse()
	return segs
}

func parseFenceOpen(line string) (marker, lang string, ok bool) {
	trimmed := strings.TrimRightFunc(line, unicode.IsSpace)
	if strings.TrimLeftFunc(trimmed, unicode.IsSpace) != trimmed {
		// Indented fences are left to Glamour.
		return "", "", false
	}
	if strings.HasPrefix(trimmed, "```") {
		n := countPrefix(trimmed, '`')
		if n < 3 {
			return "", "", false
		}
		return strings.Repeat("`", n), strings.TrimSpace(trimmed[n:]), true
	}
	if strings.HasPrefix(trimmed, "~~~") {
		n := countPrefix(trimmed, '~')
		if n < 3 {
			return "", "", false
		}
		return strings.Repeat("~", n), strings.TrimSpace(trimmed[n:]), true
	}
	return "", "", false
}

func isFenceClose(line, marker string) bool {
	trimmed := strings.TrimRightFunc(line, unicode.IsSpace)
	if strings.TrimLeftFunc(trimmed, unicode.IsSpace) != trimmed {
		return false
	}
	ch := rune(marker[0])
	n := countPrefix(trimmed, ch)
	return n >= len(marker) && strings.TrimSpace(trimmed[n:]) == ""
}

func countPrefix(s string, ch rune) int {
	n := 0
	for _, r := range s {
		if r != ch {
			break
		}
		n++
	}
	return n
}

func renderBoxedCode(code, lang string, width int, theme Theme) string {
	var buf bytes.Buffer
	chromaTheme := "dracula"
	if theme == ThemeLight {
		chromaTheme = "github"
	}
	if err := quick.Highlight(&buf, code, lang, "terminal256", chromaTheme); err != nil {
		buf.Reset()
		buf.WriteString(code)
	}
	body := strings.TrimRight(buf.String(), "\n")
	if body == "" {
		body = " "
	}

	border := "245"
	if theme == ThemeLight {
		border = "238"
	}
	inner := width - 4 // borders + horizontal padding
	if inner < 8 {
		inner = max(1, width-2)
	}

	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(border)).
		Padding(0, 1).
		Width(inner).
		Render(body)
}
