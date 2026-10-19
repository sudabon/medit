package render

import (
	"strings"
	"unicode"

	"github.com/charmbracelet/x/ansi"
)

// wrapOverflow breaks rendered lines wider than width. Prose is rendered by
// Glamour without wrapping because its wrapper only breaks at spaces and
// pushes whole runs of Japanese text onto the next line. Lines break at
// spaces and between wide (CJK) characters, honoring basic kinsoku rules.
// Continuation lines repeat the leading indent and block-quote bars.
func wrapOverflow(rendered string, width int) string {
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

type token struct {
	s string
	w int
}

func tokenize(line string) []token {
	var toks []token
	var state byte
	for len(line) > 0 {
		seq, w, n, newState := ansi.DecodeSequence(line, state, nil)
		toks = append(toks, token{seq, w})
		state = newState
		line = line[n:]
	}
	return toks
}

func wrapLine(line string, width int) []string {
	keepW, markerW := continuationIndent(ansi.Strip(line))
	if keepW+markerW >= width/2 {
		keepW, markerW = 0, 0
	}
	indent := ansi.Truncate(line, keepW, "") + "\x1b[0m" + strings.Repeat(" ", markerW)
	prefixW := keepW + markerW

	var out []string
	cur := []token{}
	curW := 0
	// cut is the index in cur before which the line may be broken; 0 = none.
	cut := 0
	minCut := 0
	var prev string

	emit := func(toks []token) {
		var b strings.Builder
		for _, t := range toks {
			b.WriteString(t.s)
		}
		s := strings.TrimRightFunc(b.String(), unicode.IsSpace)
		if activeSGR(toks) != "" {
			s += "\x1b[0m"
		}
		out = append(out, s)
	}
	startLine := func(carry []token, sgr string) {
		cur = tokenize(indent + sgr)
		curW = prefixW
		minCut = len(cur)
		cut = 0
		for len(carry) > 0 && carry[0].s == " " {
			carry = carry[1:]
		}
		for _, t := range carry {
			cur = append(cur, t)
			curW += t.w
		}
	}

	for _, t := range tokenize(line) {
		if t.w > 0 {
			if prev != "" && curW > prefixW && canBreak(prev, t.s) && len(cur) > minCut {
				cut = len(cur)
			}
			if curW+t.w > width && t.s != " " {
				head, carry := cur, []token(nil)
				if cut > minCut {
					head, carry = cur[:cut], cur[cut:]
				}
				emit(head)
				startLine(carry, activeSGR(head))
			}
			prev = t.s
		}
		cur = append(cur, t)
		curW += t.w
	}
	emit(cur)
	return out
}

// activeSGR returns the SGR sequences still in effect after toks.
func activeSGR(toks []token) string {
	var b strings.Builder
	for _, t := range toks {
		if !strings.HasPrefix(t.s, "\x1b[") || !strings.HasSuffix(t.s, "m") {
			continue
		}
		if t.s == "\x1b[0m" || t.s == "\x1b[m" {
			b.Reset()
			continue
		}
		b.WriteString(t.s)
	}
	return b.String()
}

const (
	noLineStart = "、。，．,.)）]］}｝」』】〕〉》!?！？:;：；・ーゃゅょっぁぃぅぇぉゎャュョッァィゥェォヮ々〜…‥%％"
	noLineEnd   = "（(［[｛{「『【〔〈《"
)

// canBreak reports whether a line may break between graphemes a and b.
func canBreak(a, b string) bool {
	if b == " " {
		return false
	}
	if a == " " {
		return true
	}
	if strings.Contains(noLineStart, b) || strings.Contains(noLineEnd, a) {
		return false
	}
	if a == "-" {
		return true
	}
	return ansi.StringWidth(a) > 1 || ansi.StringWidth(b) > 1
}

// continuationIndent returns the widths of the prefix to repeat on wrapped
// lines: leading spaces and "│ " quote bars kept as-is, followed by blanks
// for a list marker.
func continuationIndent(plain string) (keepW, markerW int) {
	rest := plain
	for {
		switch {
		case strings.HasPrefix(rest, " "):
			keepW++
			rest = rest[1:]
		case strings.HasPrefix(rest, "│ "):
			keepW += 2
			rest = rest[len("│ "):]
		default:
			if marker := listMarker(rest); marker != "" {
				markerW = ansi.StringWidth(marker)
			}
			return keepW, markerW
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
