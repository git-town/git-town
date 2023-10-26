package slice

// RemoveAt provides the given list with the element at the given position removed.
func RemoveAt[S ~[]C, C comparable](list S, index int) S { //nolint:ireturn
	return append(list[:index], list[index+1:]...)
}
