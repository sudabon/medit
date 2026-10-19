package render

import (
	"regexp"
	"strings"
)

type mdPart struct {
	markdown string
	isTable  bool
}

var tableDelimiter = regexp.MustCompile(`^\|?\s*:?-+:?\s*(\|\s*:?-+:?\s*)*\|?\s*$`)

// splitTables splits top-level GFM tables out of prose so they can be laid
// out to the content width while prose is wrapped separately.
func splitTables(src string) []mdPart {
	lines := strings.Split(src, "\n")
	var parts []mdPart
	start := 0
	flush := func(end int, isTable bool) {
		if end > start {
			parts = append(parts, mdPart{strings.Join(lines[start:end], "\n"), isTable})
		}
		start = end
	}
	for i := 0; i+1 < len(lines); i++ {
		if !isTableStart(lines[i], lines[i+1]) {
			continue
		}
		flush(i, false)
		end := i + 2
		for end < len(lines) && strings.TrimSpace(lines[end]) != "" {
			end++
		}
		flush(end, true)
		i = end - 1
	}
	flush(len(lines), false)
	return parts
}

func isTableStart(header, delim string) bool {
	if header == "" || header[0] == ' ' || header[0] == '\t' || !strings.Contains(header, "|") {
		return false
	}
	return strings.Contains(delim, "|") && tableDelimiter.MatchString(delim)
}
