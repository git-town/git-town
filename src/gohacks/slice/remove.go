package slice

// Remove provides the given list without the given element.
func Remove[S ~[]C, C comparable](list S, value C) S { //nolint:ireturn
	listLen := len(list)
	if listLen == 0 {
		return list
	}
	result := make([]C, 0, listLen-1)
	for _, element := range list {
		if element != value {
			result = append(result, element)
		}
	}
	return result
}
