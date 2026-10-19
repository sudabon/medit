package render

import (
	"strings"
	"unicode"

	"github.com/charmbracelet/x/ansi"
)

// hardWrapOverflow breaks rendered lines wider than width. Glamour only wraps
// at spaces, so text without spaces (Japanese, long URLs) would otherwise run
// past the content column and be cut off by the viewport.
// Continuation lines repeat the leading indent and block-quote bars.
func hardWrapOverflow(rendered string, width int) string {
	if width < 1 {
		return rendered
	}
	lines := strings.Split(rendered, "\n")
	out := make([]string, 0, len(lines))
	for _, line := range lines {
		if ansi.StringWidth(line) <= width {
			out = append(out, line)
			continue
		}
		out = append(out, wrapLine(line, width)...)
	}
	return strings.Join(out, "\n")
}

func wrapLine(line string, width int) []string {
	indent := continuationIndent(ansi.Strip(line))
	indentW := ansi.StringWidth(indent)
	if indentW >= width/2 {
		indent, indentW = "", 0
	}

	first := ansi.Truncate(line, width, "")
	rest := ansi.TruncateLeft(line, ansi.StringWidth(first), "")
	rest = strings.TrimRightFunc(rest, unicode.IsSpace)

	wrapped := []string{first}
	if ansi.StringWidth(rest) == 0 {
		return wrapped
	}
	for _, part := range strings.Split(ansi.Hardwrap(rest, width-indentW, false), "\n") {
		wrapped = append(wrapped, indent+part)
	}
	return wrapped
}

// continuationIndent returns the prefix to repeat on wrapped lines: leading
// spaces and "│ " quote bars, with list markers replaced by blanks.
func continuationIndent(plain string) string {
	var b strings.Builder
	rest := plain
	for {
		switch {
		case strings.HasPrefix(rest, " "):
			b.WriteByte(' ')
			rest = rest[1:]
		case strings.HasPrefix(rest, "│ "):
			b.WriteString("│ ")
			rest = rest[len("│ "):]
		default:
			if marker := listMarker(rest); marker != "" {
				b.WriteString(strings.Repeat(" ", ansi.StringWidth(marker)))
			}
			return b.String()
		}
	}
}

func listMarker(s string) string {
	for _, m := range []string{"• ", "[✓] ", "[ ] "} {
		if strings.HasPrefix(s, m) {
			return m
		}
	}
	i := 0
	for i < len(s) && s[i] >= '0' && s[i] <= '9' {
		i++
	}
	if i > 0 && strings.HasPrefix(s[i:], ". ") {
		return s[:i+2]
	}
	return ""
}
