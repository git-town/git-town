package stringss

import "strings"

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
