package slice

// Remove removes the given element from the given slice.
func Remove[S ~[]C, C comparable](list *S, value C) {
	listLen := len(*list)
	if listLen == 0 {
		return
	}
	result := make([]C, 0, listLen-1)
	for _, element := range *list {
		if element != value {
			result = append(result, element)
		}
	}
	*list = result
}
