package slice

// Hoist provides the given list with the given element moved to the first position.
func Hoist[S ~[]C, C comparable](list S, needle C) S { //nolint:ireturn
	result := make([]C, 0, len(list))
	hasNeedle := false
	for _, element := range list {
		if element == needle {
			hasNeedle = true
		} else {
			result = append(result, element)
		}
	}
	if hasNeedle {
		result = append([]C{needle}, result...)
	}
	return result
}
