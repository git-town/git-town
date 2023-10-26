package slice

// LowerLast provides the given slice with the last element of the given type moved to the last position in the list.
func LowerAll[C comparable](haystack []C, needle C) []C {
	result := make([]C, 0, len(haystack))
	hasNeedle := false
	for _, element := range haystack {
		if element == needle {
			hasNeedle = true
		} else {
			result = append(result, element)
		}
	}
	if hasNeedle {
		result = append(result, needle)
	}
	return result
}
