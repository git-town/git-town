package slice

// Window provides size elements surrounding the given cursor.
func Window[S ~[]C, C comparable](elements S, cursor int, size int) (window S, cursorRow int) {
	return elements, 0
}
