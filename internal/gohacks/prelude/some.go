package prelude

// Some instantiates a new Option containing the given value.
func Some[T any](value T) Option[T] {
	return Option[T]{&value}
}
