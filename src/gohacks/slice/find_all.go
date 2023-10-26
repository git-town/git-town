package slice

// FindAll provides the indexes of all occurrences of the given element in the given list.
func FindAll[C comparable](list []C, element C) []int {
	result := []int{}
	for l, li := range list {
		if li == element {
			result = append(result, l)
		}
	}
	return result
}
