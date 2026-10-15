package ui

import (
	"strings"
	"unicode"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
)

// searchMatch is one occurrence of the query in the rendered output.
// start and end are display-cell columns on the line ([start, end)).
type searchMatch struct {
	line  int
	start int
	end   int
}

var (
	matchStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("16")).
			Background(lipgloss.Color("228"))
	currentMatchStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("16")).
				Background(lipgloss.Color("208")).
				Bold(true)
)

// findMatches returns non-overlapping occurrences of query in the visible
// text of lines (ANSI sequences are ignored). Matching is smart-case:
// case-insensitive unless the query contains an uppercase letter.
func findMatches(lines []string, query string) []searchMatch {
	if query == "" {
		return nil
	}
	fold := !hasUpper(query)
	q := []rune(query)
	if fold {
		q = lowerRunes(q)
	}

	var matches []searchMatch
	for i, line := range lines {
		orig := []rune(ansi.Strip(line))
		if len(orig) < len(q) {
			continue
		}
		hay := orig
		if fold {
			hay = lowerRunes(orig)
		}
		for j := 0; j+len(q) <= len(hay); {
			if !runesEqual(hay[j:j+len(q)], q) {
				j++
				continue
			}
			start := ansi.StringWidth(string(orig[:j]))
			end := start + ansi.StringWidth(string(orig[j:j+len(q)]))
			matches = append(matches, searchMatch{line: i, start: start, end: end})
			j += len(q)
		}
	}
	return matches
}

// highlightLines returns a copy of lines with every match highlighted and
// matches[current] drawn with the current-match style.
func highlightLines(lines []string, matches []searchMatch, current int) []string {
	out := append([]string(nil), lines...)
	for i := 0; i < len(matches); {
		j := i
		for j < len(matches) && matches[j].line == matches[i].line {
			j++
		}
		out[matches[i].line] = highlightLine(lines[matches[i].line], matches[i:j], current-i)
		i = j
	}
	return out
}

// highlightLine wraps each match on a single line. current is the index into
// ms of the current match, or out of range when it is on another line.
func highlightLine(line string, ms []searchMatch, current int) string {
	var b strings.Builder
	pos := 0
	for k, m := range ms {
		b.WriteString(ansi.Cut(line, pos, m.start))
		style := matchStyle
		if k == current {
			style = currentMatchStyle
		}
		b.WriteString(style.Render(ansi.Strip(ansi.Cut(line, m.start, m.end))))
		pos = m.end
	}
	b.WriteString(ansi.TruncateLeft(line, pos, ""))
	return b.String()
}

func hasUpper(s string) bool {
	for _, r := range s {
		if unicode.IsUpper(r) {
			return true
		}
	}
	return false
}

func lowerRunes(rs []rune) []rune {
	out := make([]rune, len(rs))
	for i, r := range rs {
		out[i] = unicode.ToLower(r)
	}
	return out
}

func runesEqual(a, b []rune) bool {
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
