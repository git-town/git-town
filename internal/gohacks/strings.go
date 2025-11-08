package gohacks

import (
	"strings"

	"github.com/git-town/git-town/v22/internal/gohacks/stringslice"
)

func EscapeNewLines(text string) string {
	return strings.ReplaceAll(text, "\n", "\\n")
}

// LeadingWhitespace provides the leading whitespace in the given string.
func LeadingWhitespace(line string) string {
	result := strings.Builder{}
	for _, r := range line {
		if r == ' ' || r == '\t' {
			result.WriteRune(r)
		} else {
			break
		}
	}
	return result.String()
}

// Indents provides the given text where each line is indented by the given amount of spaces.
func IndentLines(text string, amount int) string {
	lines := strings.Split(text, "\n")
	lines = stringslice.Indent(lines, strings.Repeat(" ", amount))
	return strings.Join(lines, "\n")
}
