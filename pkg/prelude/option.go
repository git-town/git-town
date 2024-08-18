package prelude

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// Option provides infrastructure for optional (nullable) values
// that is fully enforced by the type checker.
// Matching the data architecture of this codebase, this Option
// provides copies of the optional value, i.e. works only for const and copyable values.
// If you need direct access to the optional value, i.e. don't want a copy, use an OptionP instead.
// The zero value is the None option.
//
// Option is worth the overhead because it removes one of the many possible meanings (optionality)
// from pointer values. This means a pointer in this codebase implies mutability and nothing else.
//
// Compare Options using their .Equal method since direct comparison using == doesn't work properly.
type Option[T any] struct {
	value *T
}

// indicates whether the given other Option has the same value as this Option
func (self Option[T]) Equal(other Option[T]) bool {
	selfValue, hasSelfValue := self.Get()
	otherValue, hasOtherValue := other.Get()
	if !hasSelfValue && !hasOtherValue {
		return true
	}
	if hasSelfValue != hasOtherValue {
		return false
	}
	return reflect.DeepEqual(selfValue, otherValue)
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

// MarshalJSON is used when serializing this Option to JSON.
func (self Option[T]) MarshalJSON() ([]byte, error) {
	if value, hasValue := self.Get(); hasValue {
		return json.Marshal(value)
	}
	return json.Marshal(nil)
}

// Performs a logical OR operation on this option and the given option:
// Returns this option if it is some, otherwise the given option.
func (self Option[T]) Or(other Option[T]) Option[T] {
	if self.IsSome() {
		return self
	}
	return other
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

// UnmarshalJSON is used when de-serializing JSON into an Option.
func (self *Option[T]) UnmarshalJSON(b []byte) error {
	if string(b) == "null" {
		self.value = nil
		return nil
	}
	var value T
	self.value = &value
	return json.Unmarshal(b, &self.value)
}

// Creates a new Option containing None if the given value is the zero value, otherwise Some.
func NewOption[T comparable](value T) Option[T] {
	var zero T
	if value == zero {
		return None[T]()
	}
	return Some(value)
}
