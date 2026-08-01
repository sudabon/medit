package render

// PreferredContentWidth is the default column width for readable Markdown.
const PreferredContentWidth = 88

// ContentWidth returns the fixed content column width for a terminal width.
// It prefers PreferredContentWidth, shrinking on narrow terminals while
// leaving a little side padding when possible.
func ContentWidth(termWidth int) int {
	if termWidth < 1 {
		return 1
	}
	maxContent := termWidth - 4
	if maxContent < 40 {
		return termWidth
	}
	if PreferredContentWidth <= maxContent {
		return PreferredContentWidth
	}
	return maxContent
}

// LeftPad returns spaces to prepend so a content column is horizontally centered.
func LeftPad(termWidth, contentWidth int) int {
	if contentWidth >= termWidth || termWidth < 1 {
		return 0
	}
	return (termWidth - contentWidth) / 2
}
