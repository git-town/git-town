package slice

// RemoveAt removes the element at the given position from the given list.
func RemoveAt[S ~[]C, C comparable](list S, index int) S {
	return append((list)[:index], (list)[index+1:]...)
}
