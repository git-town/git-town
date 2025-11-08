package stringslice

import "strings"

// IndentNonEmpty prepends the given indentation to each non-empty line
func IndentNonEmpty(lines []string, indentation string) []string {
	result := make([]string, len(lines))
	for i, line := range lines {
		if strings.TrimSpace(line) != "" {
			result[i] = indentation + strings.TrimLeft(line, " \t")
		} else {
			result[i] = line
		}
	}
	return result
}
