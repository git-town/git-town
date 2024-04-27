package prelude

// Some instantiates a new Option containing the given value.
func SomeP[T any](value T) OptionP[T] {
	return OptionP[T]{&value}
}
