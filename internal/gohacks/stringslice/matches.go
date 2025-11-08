package stringslice

import "strings"

// Matches indicates whether the given string slices match.
func Matches(lines1, lines2 []string) bool {
	if len(lines1) != len(lines2) {
		return false
	}
	for l, line2 := range lines2 {
		if strings.TrimSpace(lines1[l]) != strings.TrimSpace(line2) {
			return false
		}
	}
	return true
}
