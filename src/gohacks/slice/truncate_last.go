package slice

// TruncateLast provides the given list without its last element.
func TruncateLast[S ~[]C, C comparable](list S) S { //nolint: ireturn
	listLength := len(list)
	if listLength == 0 {
		return list
	}
	return list[:listLength-1]
}
