package prelude

// SomeP instantiates a new OptionP containing the given value.
// The value must exist, i.e. the pointer must not be nil.
func SomeP[T any](value *T) OptionP[T] {
	if value == nil {
		panic("You gave nil to SomeP")
	}
	return OptionP[T]{value}
}
