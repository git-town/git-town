package stringslice

import "strings"

func LinesWithPrefix(lines []string, prefix string) []string {
	result := make([]string, 0, 1)
	for _, line := range lines {
		if strings.HasPrefix(line, prefix) {
			result = append(result, line)
		}
	}
	return result
}
