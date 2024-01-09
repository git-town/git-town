package slice

import "slices"

// FindMany provides the indixes of all needles in the give haystack.
func FindMany[S ~[]C, C comparable](haystack S, needles S) []int {
	result := make([]int, 0, len(needles))
	for _, needle := range needles {
		pos := slices.Index(haystack, needle)
		if pos != -1 {
			result = append(result, pos)
		}
	}
	return result
}
