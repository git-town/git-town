package prelude

// MutableNone instantiates an empty MutableOption.
func MutableNone[T any]() OptionalMutable[T] {
	return OptionalMutable[T]{nil}
}
