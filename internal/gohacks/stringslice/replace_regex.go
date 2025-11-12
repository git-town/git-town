package stringslice

import "regexp"

// ReplaceRegex replaces all matches of the given regex with the given replacement in the given strings.
func ReplaceRegex(texts []string, regex *regexp.Regexp, replacement string) []string {
	result := make([]string, len(texts))
	for t, text := range texts {
		result[t] = regex.ReplaceAllString(text, replacement)
	}
	return result
}
