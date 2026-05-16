package stringss

import "strings"

// IndentLines provides the given text where each line is indented by the given amount of spaces.
func IndentLines(text string, amount int) string {
	lines := strings.Split(text, "\n")
	result := make([]string, len(lines))
	indent := strings.Repeat(" ", amount)
	for l, line := range lines {
		if len(line) > 0 {
			result[l] = indent + line
		}
	}
	return strings.Join(result, "\n")
}
