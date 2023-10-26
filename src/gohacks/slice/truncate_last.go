package slice

// TruncateLast provides the given list with its last element removed.
func TruncateLast[S ~[]C, C comparable](list S) S { //nolint:ireturn
	listLength := len(list)
	if listLength == 0 {
		return list
	}
	return list[:listLength-1]
}
