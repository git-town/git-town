package gohacks

import "strings"

func EscapeNewLines(text string) string {
	return strings.ReplaceAll(text, "\n", "\\n")
}

func IndentLines(text string, indent int) string {
	result := strings.Builder{}
	for _, line := range strings.Split(text, "\n") {
		for range indent {
			result.WriteString(" ")
		}
		result.WriteString(line)
		result.WriteString("\n")
	}
	return result.String()
}
