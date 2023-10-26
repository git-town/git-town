package slice

// LowerLast moves all occurrences of the given element in the given list to the end of the list.
func LowerAll[S ~[]C, C comparable](haystack *S, needle C) {
	result := make([]C, 0, len(*haystack))
	hasNeedle := false
	for _, element := range *haystack {
		if element == needle {
			hasNeedle = true
		} else {
			result = append(result, element)
		}
	}
	if hasNeedle {
		result = append(result, needle)
	}
	*haystack = result
}
