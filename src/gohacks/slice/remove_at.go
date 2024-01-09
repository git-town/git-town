package slice

// RemoveAt removes the element at the given position from the given list.
func RemoveAt[S ~[]C, C comparable](list S, index int) S { //nolint:ireturn
	return append((list)[:index], (list)[index+1:]...)
}
