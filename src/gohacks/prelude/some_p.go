package prelude

// Some instantiates a new Option containing the given value.
func SomeP[T any](value *T) OptionP[T] {
	if value == nil {
		panic("You gave nil to SomeP")
	}
	return OptionP[T]{value}
}
