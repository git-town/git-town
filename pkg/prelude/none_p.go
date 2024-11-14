package prelude

// NoneP instantiates an OptionP of the given type containing None.
func NoneP[T any]() MutableOption[T] {
	return MutableOption[T]{nil}
}
