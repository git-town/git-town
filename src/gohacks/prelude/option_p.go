package prelude

import (
	"fmt"
)

// OptionP is an Option that provides direct access to the encapsulated value using pointers.
type OptionP[T any] struct {
	value *T
}

// Get provides a copy of the contained value
// as well as an indicator whether that value exists.
func (self OptionP[T]) Get() (value *T, hasValue bool) {
	if self.IsSome() {
		return self.value, true
	}
	var empty T
	return &empty, false
}

// GetOrDefault provides a copy of the contained value.
// If this option contains nothing, you get the zero value of the contained type.
func (self OptionP[T]) GetOrDefault() *T {
	if value, has := self.Get(); has {
		return value
	}
	var empty T
	return &empty
}

// GetOrElse provides a copy of the contained value.
// If this option contains nothing, you get a copy of the given alternative value.
func (self OptionP[T]) GetOrElse(other *T) *T {
	if value, has := self.Get(); has {
		return value
	}
	return other
}

// GetOrPanic provides a copy of the contained value.
// If this option nothing, this method panics.
func (self OptionP[T]) GetOrPanic() *T {
	if value, has := self.Get(); has {
		return value
	}
	panic("value not present")
}

// GetOrPanic provides direct access to the contained value via a pointer.
// If this option nothing, this method panics.
func (self OptionP[T]) GetPOrPanic() *T {
	if self.IsSome() {
		return self.value
	}
	panic("value not present")
}

// IsNone indicates whether this option instance contains nothing.
func (self OptionP[T]) IsNone() bool {
	return self.value == nil
}

// IsSome indicates whether this option instance contains a value.
func (self OptionP[T]) IsSome() bool {
	return self.value != nil
}

// String provides the string serialization of the contained value.
// If this option contains nothing, you get an empty string.
func (self OptionP[T]) String() string {
	return self.StringOr("")
}

// StringOr provideds the string serialization of the contained value.
// If this option contains nothing, you get the given alternative string representation.
func (self OptionP[T]) StringOr(other string) string {
	if self.IsNone() {
		return other
	}
	return fmt.Sprint(self.value)
}
