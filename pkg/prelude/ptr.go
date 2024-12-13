package prelude

func Ptr[A any](value A) *A {
	return &value
}
