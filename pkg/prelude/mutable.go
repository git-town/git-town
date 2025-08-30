package prelude

// Mutable respresents a mutable value.
// A Mutable is always mutable, even if passed by value.
type Mutable[T any] struct {
	// the enclosed mutable value, okay to mutate it directly
	Value *T
}

// provides an non-mutable copy of the contained mutable value
func (self Mutable[T]) Immutable() T { //nolint:ireturn
	return *self.Value
}

// MutableNone instantiates an empty MutableOption.
func MutableNone[T any]() OptionalMutable[T] {
	return OptionalMutable[T]{nil}
}

// MutableSome instantiates a new OptionP containing the given value.
// The value must exist, i.e. the pointer must not be nil.
func MutableSome[T any](value *T) OptionalMutable[T] {
	if value == nil {
		panic("Cannot create a SomeP out of a nil pointer")
	}
	return OptionalMutable[T]{value}
}

func NewMutable[T any](value *T) Mutable[T] {
	return Mutable[T]{value}
}
