package prelude

// Mutable respresents a mutable value.
// A Mutable is always mutable, even if passed by value.
type Mutable[T any] struct {
	// the enclosed mutable value, okay to mutate it directly
	Value *T
}

func NewMutable[T any](value *T) Mutable[T] {
	return Mutable[T]{value}
}

// provides an non-mutable copy of the contained mutable value
func (self Mutable[T]) Copy() T { //nolint:ireturn
	return *self.Value
}
