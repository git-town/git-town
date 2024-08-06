package slice

import "slices"

// indicates whether the given haystack contains any of the given needles
func ContainsAny[C comparable](haystack []C, needles []C) bool {
	for _, needle := range needles {
		if slices.Contains(haystack, needle) {
			return true
		}
	}
	return false
}
