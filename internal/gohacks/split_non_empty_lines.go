package gohacks

import "strings"

// splitNonEmptyLines splits the input by newlines and returns non-empty, trimmed lines.
func SplitNonEmptyLines(output string) []string {
	if output == "" {
		return []string{}
	}
	lines := strings.Split(output, "\n")
	result := make([]string, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			result = append(result, line)
		}
	}
	return result
}
