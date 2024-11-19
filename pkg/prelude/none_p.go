package prelude

// MutableNone instantiates a non-existing MutableOption.
func MutableNone[T any]() OptionalMut[T] {
	return OptionalMut[T]{nil}
}
