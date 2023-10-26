package slice

// Hoist provides the given slice with the given element moved to the first position.
func Hoist[S ~[]C, C comparable](list S, element C) S { //nolint:ireturn
	result := make([]C, 0, len(list))
	hasElement := false
	for l := range list {
		if list[l] == element {
			hasElement = true
		} else {
			result = append(result, list[l])
		}
	}
	if hasElement {
		result = append([]C{element}, result...)
	}
	return result
}
