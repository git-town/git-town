package prelude

// Ptr provides a pointer to the given value.
//
// This is useful for getting pointers to literals of basic values.
// This is not possible in Go:
// x := &"hello"
// This function makes it possible:
// x := Ptr("hello")
func Ptr[T any](value T) *T {
	return &value
}
