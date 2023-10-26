package slice

// TruncateLast removes the last element from the given list.
func TruncateLast[S ~[]C, C comparable](list *S) {
	listLength := len(*list)
	if listLength == 0 {
		return
	}
	*list = (*list)[:listLength-1]
}
