package prelude

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/git-town/git-town/v22/pkg/equal"
)

// Option encodes invariants around optional (nullable) immutable values in the type system.
// The zero value is the None option.
//
// For optional values that should not be copied, please use OptionalMutable.
// Compare Options using their `Equal` or `EqualSome` methods,
// direct comparison using == doesn't work.
//
// Why pointers are not a good solution for optional values:
//
//  1. A pointer has many meanings: optional, mutable, too large to pass by value.
//     The Option type documents that a value is optional (and immutable, or mutable when using OptionalMutable).
//
//  2. A pointer looks the same before and after you checked it for nil.
//     Pointers therefore carry the risk of being checked too often or too little,
//     leading to unnecessary boilerplate or bugs.
//     Options get checked exactly once, leading to the least amount of boilerplate code.
//
//  3. Less "!= nil" in your code.
type Option[T any] struct {
	value *T
}

// Equal indicates whether the given other Option has the same value as this Option
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

// EqualSome indicates whether this option contains the given value
func (self Option[T]) EqualSome(other T) bool {
	if value, hasValue := self.Get(); hasValue {
		return reflect.DeepEqual(value, other)
	}
	return false
}

// Get provides a copy of the contained value
// as well as an indicator whether that value exists.
func (self Option[T]) Get() (value T, hasValue bool) { //nolint:ireturn,nonamedreturns
	if self.IsSome() {
		return *self.value, true
	}
	var empty T
	return empty, false
}

// GetOr provides a copy of the contained value.
// If this option contains nothing, you get a copy of the given alternative value.
func (self Option[T]) GetOr(other T) T { //nolint:ireturn
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

// GetOrZero provides a copy of the contained value.
// If this option contains nothing, you get the zero value of the contained type.
func (self Option[T]) GetOrZero() T { //nolint:ireturn
	if value, has := self.Get(); has {
		return value
	}
	var empty T
	return empty
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

// Or performs a logical OR operation on this option and the given option:
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
	if value, has := self.Get(); has {
		return fmt.Sprint("Some(", value, ")")
	}
	return "None"
}

// StringOr provides the string serialization of the contained value.
// If this option contains nothing, you get the given alternative string representation.
func (self Option[T]) StringOr(other string) string {
	if value, has := self.Get(); has {
		return fmt.Sprint(value)
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

// NewOption creates a new Option containing None if the given value is the zero value, otherwise Some.
func NewOption[T any](value T) Option[T] {
	var zero T
	if equal.Equal(value, zero) {
		return None[T]()
	}
	return Some(value)
}

// None instantiates an empty Option of the given type.
func None[T any]() Option[T] {
	return Option[T]{nil}
}

// Some instantiates a new Option containing the given value.
func Some[T any](value T) Option[T] {
	return Option[T]{&value}
}
