package prelude

// Mutable respresents a mutable value.
// A value wrapped in Mutable is always mutable, even if passed by reference.
type Mutable[T any] struct {
	Value *T
}

func NewMutable[T any]() Mutable[T] {
	var value T
	return Mutable[T]{&value}
}

// provides an immutable copy of the current value
func (self Mutable[T]) Get() T { //nolint:ireturn
	return *self.Value
}
