package prelude

import (
	"fmt"
)

// OptionP ("option for a pointer value") is an Option
// that provides direct access to the encapsulated value
// by storing and providing a pointer to the value.
// This is useful for mutable or singleton values,
// or values that are too large to copy around all the time.
type OptionP[T any] struct {
	value *T
}

// Get provides a copy of the contained value
// as well as an indicator whether that value exists.
func (self OptionP[T]) Get() (value *T, hasValue bool) {
	if self.IsSome() {
		return self.value, true
	}
	return nil, false
}

// GetOrPanic provides a copy of the contained value.
// Panics if this option contains nothing.
func (self OptionP[T]) GetOrPanic() *T {
	if value, has := self.Get(); has {
		return value
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
	if self.IsSome() {
		return fmt.Sprint(self.value)
	}
	return other
}
