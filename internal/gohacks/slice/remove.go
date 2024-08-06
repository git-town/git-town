package slice

// Remove provides the given list without the given element.
func Remove[S ~[]C, C comparable](list S, value C) S { //nolint:ireturn
	if len(list) == 0 {
		return list
	}
	result := make([]C, 0, len(list)-1)
	for _, element := range list {
		if element != value {
			result = append(result, element)
		}
	}
	return result
}
