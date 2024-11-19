package prelude

// MutableSome instantiates a new OptionP containing the given value.
// The value must exist, i.e. the pointer must not be nil.
func MutableSome[T any](value *T) OptionalMutable[T] {
	if value == nil {
		panic("Cannot create a SomeP out of a nil pointer")
	}
	return OptionalMutable[T]{value}
}
