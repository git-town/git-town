package stringslice

import "strings"

// EqualIgnoreWhitespace indicates whether the given string slices match,
// ignoring whitespace before and after each string.
func EqualIgnoreWhitespace(strings1, strings2 []string) bool {
	if len(strings1) != len(strings2) {
		return false
	}
	for s, string1 := range strings1 {
		if strings.TrimSpace(strings2[s]) != strings.TrimSpace(string1) {
			return false
		}
	}
	return true
}
