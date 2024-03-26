package util

import (
	"strings"
)

// LastNLines returns the last n lines from a given string.
func LastNLines(input string, n int) string {
	lines := strings.Split(input, "\n")
	lineCount := len(lines)

	// Adjust if requesting more lines than are available.
	if n > lineCount {
		n = lineCount
	}

	lastLines := lines[lineCount-n : lineCount]
	return strings.Join(lastLines, "\n")
}
