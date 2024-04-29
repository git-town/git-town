package prelude

// SomeP instantiates a new OptionP containing the given value.
// The value must exist, i.e. the pointer must not be nil.
func SomeP[T any](value *T) OptionP[T] {
	if value == nil {
		panic("Cannot create a SomeP out of a nil pointer")
	}
	return OptionP[T]{value}
}
