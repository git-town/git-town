package prelude

import (
	"fmt"
)

// Option provides infrastructure for nullable values that is enforced by the type checker.
// Since all types used in Git Town implement the fmt.Stringer interface,
// we can narrow the allowed types to fmt.Stringer.
// The zero value is the None option.
//
// We tried using pointers to express optionality before but it doesn't work well.
// There are too many situation where a pointer expression happily passes the type checker
// and then panics at runtime.
// Go sometimes de-references pointers and sometimes it doesn't.
// Pointers have too many meanings: reference, mutability, poor-man optionality.
// Better to have a dedicated facility for optionality and only that.
type Option[T any] struct {
	Value *T
}

// Get provides the contained value as well as an indicator whether that value exists.
func (self Option[T]) Get() (value T, hasValue bool) { //nolint:ireturn
	if self.IsSome() {
		return *self.Value, true
	}
	var empty T
	return empty, false
}

// Get provides the contained value as well as an indicator whether that value exists.
func (self Option[T]) GetP() (value *T, hasValue bool) { //nolint:ireturn
	if self.IsSome() {
		return self.Value, true
	}
	var empty T
	return &empty, false
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

// GetOrPanic provides the contained value. If this option nothing,
// this method panics.
func (self Option[T]) GetPOrPanic() *T { //nolint:ireturn
	if self.IsSome() {
		return self.Value
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
	if self.IsNone() {
		return other
	}
	return fmt.Sprint(self.Value)
}
