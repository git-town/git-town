package gohacks

import "fmt"

// Option implements type-safe handling of optional values.
// Since all values in Git Town implement the fmt.Stringer interface,
// we can narrow the allowed types to it.
//
// Go doesn't provide good handling of optional values out of the box.
// Using pointers to indicate optionality results in many null-pointer-exceptions at runtime.
type Option[T fmt.Stringer] struct {
	Value *T
}

// NewOption instantiates a new option containing the given value.
func NewOption[T fmt.Stringer](value T) Option[T] {
	return Option[T]{&value}
}

// NewOptionNone instantiates a new option containing nothing.
func NewOptionNone[T fmt.Stringer]() Option[T] {
	return Option[T]{nil}
}

// Get provides the contained value as well as an indicator whether that value exists.
func (self Option[T]) Get() (value T, hasValue bool) { //nolint:ireturn
	if self.IsSome() {
		return *self.Value, true
	}
	var empty T
	return empty, false
}

// GetOrDefault provides the contained value. If this option contains nothing,
// you get the zero value of the contained type.
func (self Option[T]) GetOrDefault() T { //nolint:ireturn
	if value, has := self.Get(); has {
		return value
	}
	var empty T
	return empty
}

// GetOrElse provides the contained value. If this option contains nothing,
// you get the given alternative value.
func (self Option[T]) GetOrElse(other T) T { //nolint:ireturn
	if value, has := self.Get(); has {
		return value
	}
	return other
}

// GetOrPanic provides the contained value. If this option nothing,
// this method panics.
func (self Option[T]) GetOrPanic() T { //nolint:ireturn
	if value, has := self.Get(); has {
		return value
	}
	panic("value not present")
}

// IsNone indicates whether this option instance contains nothing.
func (self Option[T]) IsNone() bool {
	return self.Value == nil
}

// IsSome indicates whether this option instance contains a value.
func (self Option[T]) IsSome() bool {
	return self.Value != nil
}

// String provides the string serialization of the contained value.
// If this option contains nothing, you get an empty string.
func (self Option[T]) String() string {
	return self.StringOr("")
}

// StringOr provideds the string serialization of the contained value.
// If this option contains nothing, you get the given alternative string representation.
func (self Option[T]) StringOr(other string) string {
	if value, has := self.Get(); has {
		return value.String()
	}
	return other
}
