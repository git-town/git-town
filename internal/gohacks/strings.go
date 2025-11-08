package gohacks

import "strings"

func EscapeNewLines(text string) string {
	return strings.ReplaceAll(text, "\n", "\\n")
}

// IndentLines provides the given text where each line is indented by the given amount of spaces.
func IndentLines(text string, amount int) string {
	lines := strings.Split(text, "\n")
	result := make([]string, len(lines))
	for l, line := range lines {
		result[l] = strings.Repeat(" ", amount) + line
	}
	return strings.Join(result, "\n")
}
