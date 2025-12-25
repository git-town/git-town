package slice

import "slices"

// Remove provides the given list without the given element.
func Remove[S ~[]C, C comparable](haystack S, needles ...C) S { //nolint:ireturn
	if len(haystack) == 0 {
		return haystack
	}
	result := make([]C, 0, len(haystack)-len(needles))
	for _, element := range haystack {
		if !slices.Contains(needles, element) {
			result = append(result, element)
		}
	}
	return result
}
