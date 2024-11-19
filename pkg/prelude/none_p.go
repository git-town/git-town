package prelude

// MutableNone instantiates a non-existing MutableOption.
func MutableNone[T any]() MutableOption[T] {
	return MutableOption[T]{nil}
}
