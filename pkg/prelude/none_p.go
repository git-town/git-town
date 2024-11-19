package prelude

// MutableNone instantiates an empty MutableOption.
func MutableNone[T any]() OptionalMut[T] {
	return OptionalMut[T]{nil}
}
