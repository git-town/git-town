package prelude

// Mutable respresents a mutable value.
// A Mutable always correctly mutates its encapsulated value, even if the Mutable gets copied or passed by reference.
type Mutable[T any] struct {
	// the enclosed mutable value, feel free to mutate it directly
	Value *T
}

func NewMutable[T any](value *T) Mutable[T] {
	return Mutable[T]{value}
}

// provides an immutable copy of the enclosed mutable value
func (self Mutable[T]) Get() T { //nolint:ireturn
	return *self.Value
}
