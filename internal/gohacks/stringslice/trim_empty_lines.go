package stringslice

import "strings"

// TrimEmptyLines removes leading and trailing empty lines from the given lines.
func TrimEmptyLines(lines []string) []string {
	result := make([]string, 0, len(lines))
	for _, line := range lines {
		if len(result) > 0 || strings.TrimSpace(line) != "" {
			result = append(result, line)
		}
	}
	// Trim trailing empty lines
	for len(result) > 0 && strings.TrimSpace(result[len(result)-1]) == "" {
		result = result[:len(result)-1]
	}
	return result
}
