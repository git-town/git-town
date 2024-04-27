package prelude

// None instantiates an Option of the given type containing nothing.
func NoneP[T any]() OptionP[T] {
	return OptionP[T]{nil}
}
