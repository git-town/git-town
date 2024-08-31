package slice

// Remove provides the given list without the given element.
func Remove[S ~[]C, C comparable](haystack S, needle C) S { //nolint:ireturn
	if len(haystack) == 0 {
		return haystack
	}
	result := make([]C, 0, len(haystack)-1)
	for _, element := range haystack {
		if element != needle {
			result = append(result, element)
		}
	}
	return result
}
