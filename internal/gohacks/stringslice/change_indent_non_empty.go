package stringslice

import "strings"

// ChangeIndentNonEmpty changes the indentation of each non-empty line to the given string.
func ChangeIndentNonEmpty(lines []string, indentation string) []string {
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
