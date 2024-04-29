package prelude

import (
	"fmt"
)

// Option provides infrastructure for optional (nullable) values that is fully enforced by the type checker.
// Option provides copies of the optional value.
// If you need direct access to the optional value, use an OptionP instead.
// The zero value is the None option.
//
// A simpler approach to express optionality would be using pointers to the optional values
// since pointers can be nil.
// This is a poor-man's approach since pointers mix reference, mutability, and optionality.
// Go sometimes derefences pointers, resulting in too many situations
// where calling methods on a nil pointer happily passes the type checker
// and then panics at runtime.
// Better to have a dedicated facility just for optionality.
type Option[T any] struct {
	value *T
}

// Get provides a copy of the contained value
// as well as an indicator whether that value exists.
func (self Option[T]) Get() (value T, hasValue bool) { //nolint:ireturn
	if self.IsSome() {
		return *self.value, true
	}
	var empty T
	return empty, false
}

// GetOrDefault provides a copy of the contained value.
// If this option contains nothing, you get the zero value of the contained type.
func (self Option[T]) GetOrDefault() T { //nolint:ireturn
	if value, has := self.Get(); has {
		return value
	}
	var empty T
	return empty
}

// GetOrElse provides a copy of the contained value.
// If this option contains nothing, you get a copy of the given alternative value.
func (self Option[T]) GetOrElse(other T) T { //nolint:ireturn
	if value, has := self.Get(); has {
		return value
	}
	return other
}

// GetOrPanic provides a copy of the contained value.
// Panics if this option contains nothing.
func (self Option[T]) GetOrPanic() T { //nolint:ireturn
	if value, has := self.Get(); has {
		return value
	}
	panic("value not present")
}

// IsNone indicates whether this option instance contains nothing.
func (self Option[T]) IsNone() bool {
	return self.value == nil
}

// IsSome indicates whether this option instance contains a value.
func (self Option[T]) IsSome() bool {
	return self.value != nil
}

// String provides the string serialization of the contained value.
// If this option contains nothing, you get an empty string.
func (self Option[T]) String() string {
	return self.StringOr("")
}

// StringOr provideds the string serialization of the contained value.
// If this option contains nothing, you get the given alternative string representation.
func (self Option[T]) StringOr(other string) string {
	if self.IsSome() {
		return fmt.Sprint(self.value)
	}
	return other
}
