package render

import "strings"

// SourceLineCount returns the number of lines in source text.
// An empty string is treated as 0 lines.
func SourceLineCount(content string) int {
	if content == "" {
		return 0
	}
	return strings.Count(content, "\n") + 1
}

// ApproximateSourceLine maps a scroll fraction (0–1) to a 1-based source line
// for `vi +N`. Values are clamped to [1, sourceLines]. Empty documents return 1.
func ApproximateSourceLine(scrollPercent float64, sourceLines int) int {
	if sourceLines <= 0 {
		return 1
	}
	if scrollPercent < 0 {
		scrollPercent = 0
	}
	if scrollPercent > 1 {
		scrollPercent = 1
	}

	n := int(scrollPercent*float64(sourceLines) + 0.5)
	if n < 1 {
		n = 1
	}
	if n > sourceLines {
		n = sourceLines
	}
	return n
}
