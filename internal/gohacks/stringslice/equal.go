package stringslice

import "strings"

// EqualIgnoreWhitespace indicates whether the given string slices match,
// ignoring whitespace before and after each string.
func EqualIgnoreWhitespace(strings1, strings2 []string) bool {
	if len(strings1) != len(strings2) {
		return false
	}
	for s, string2 := range strings2 {
		if strings.TrimSpace(strings1[s]) != strings.TrimSpace(string2) {
			return false
		}
	}
	return true
}
