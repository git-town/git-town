package prelude

// Mutable respresents a mutable value.
// A Mutable is always mutable, even if passed by value.
type Mutable[T any] struct {
	// the enclosed mutable value, okay to mutate it directly
	value       *T
	initialized bool
}

// provides an non-mutable copy of the contained mutable value
func (self Mutable[T]) Immutable() T { //nolint:ireturn
	self.verify()
	return *self.value
}

// provides an non-mutable copy of the contained mutable value
func (self Mutable[T]) Value() *T {
	self.verify()
	return self.value
}

func (self Mutable[T]) verify() {
	if !self.initialized {
		panic("Found a mutable that wasn't created by the NewMutable constructor function.")
	}
}

func NewMutable[T any](value *T) Mutable[T] {
	return Mutable[T]{value, true}
}
