package slice

// Remove returns a new slice which is the given slice with the given element removed.
func Remove[S ~[]C, C comparable](list S, value C) S { //nolint:ireturn
	listLen := len(list)
	if listLen == 0 {
		return list
	}
	result := make([]C, 0, listLen-1)
	for l := range list {
		if list[l] != value {
			result = append(result, list[l])
		}
	}
	return result
}
