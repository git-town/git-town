package prelude

// None instantiates an empty Option of the given type.
func None[T any]() Option[T] {
	return Option[T]{nil}
}
