package stringslice

import "strings"

// Equal indicates whether the given string slices match.
func Equal(strings1, strings2 []string) bool {
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
