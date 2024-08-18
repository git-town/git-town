package prelude

// NoneP instantiates an OptionP of the given type containing None.
func NoneP[T any]() OptionP[T] {
	return OptionP[T]{nil}
}
