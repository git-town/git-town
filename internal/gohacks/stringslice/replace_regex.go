package stringslice

import "regexp"

func ReplaceRegex(lines []string, regex *regexp.Regexp, replacement string) []string {
	result := make([]string, len(lines))
	for i, line := range lines {
		result[i] = regex.ReplaceAllString(line, replacement)
	}
	return result
}
