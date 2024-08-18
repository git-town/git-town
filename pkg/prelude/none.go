package prelude

// None instantiates an Option of the given type containing nothing.
func None[T any]() Option[T] {
	return Option[T]{nil}
}
