package slice

// Hoist moves the given element in the given slice to the first position.
func Hoist[S ~[]C, C comparable](list *S, needle C) {
	result := make([]C, 0, len(*list))
	hasNeedle := false
	for _, element := range *list {
		if element == needle {
			hasNeedle = true
		} else {
			result = append(result, element)
		}
	}
	if hasNeedle {
		result = append([]C{needle}, result...)
	}
	*list = result
}
