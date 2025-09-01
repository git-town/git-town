package prelude

import (
	"encoding/json"
	"fmt"
)

// OptionalMutable represents a value that is both optional and mutable.
type OptionalMutable[T any] struct {
	Value       *T
	initialized bool
}

// Get provides null-safe mutable access to the contained value.
func (self OptionalMutable[T]) Get() (value *T, hasValue bool) {
	self.verify()
	if self.IsSome() {
		return self.Value, true
	}
	return nil, false
}

// GetOrPanic provides unsafe mutable access to the contained value.
func (self OptionalMutable[T]) GetOrPanic() *T {
	self.verify()
	if value, has := self.Get(); has {
		return value
	}
	panic("value not present")
}

// IsNone indicates whether this option instance contains nothing.
func (self OptionalMutable[T]) IsNone() bool {
	self.verify()
	return self.Value == nil
}

// IsSome indicates whether this option instance contains a value.
func (self OptionalMutable[T]) IsSome() bool {
	self.verify()
	return self.Value != nil
}

// MarshalJSON is used when serializing this type to JSON.
func (self OptionalMutable[T]) MarshalJSON() ([]byte, error) {
	self.verify()
	if value, hasValue := self.Get(); hasValue {
		return json.Marshal(*value)
	}
	return json.Marshal(nil)
}

// String provides the string serialization of the contained value.
// None gets serialized into an empty string.
func (self OptionalMutable[T]) String() string {
	self.verify()
	return self.StringOr("")
}

// StringOr provideds the string serialization of the contained value.
// None gets serialized into the given alternative string representation.
func (self OptionalMutable[T]) StringOr(other string) string {
	self.verify()
	if self.IsSome() {
		return fmt.Sprint(self.Value)
	}
	return other
}

// ToOption provides an immutable copy of this OptionalMut.
func (self OptionalMutable[T]) ToOption() Option[T] {
	self.verify()
	if value, hasValue := self.Get(); hasValue {
		return Some(*value)
	}
	return None[T]()
}

// UnmarshalJSON is used when de-serializing this type from JSON.
func (self *OptionalMutable[T]) UnmarshalJSON(b []byte) error {
	self.verify()
	if string(b) == "null" {
		self.Value = nil
		return nil
	}
	var value T
	self.Value = &value
	return json.Unmarshal(b, &self.Value)
}

func (self *OptionalMutable[T]) verify() {
	if !self.initialized {
		panic("Found an OptionalMutable that wasn't created by the MutableNone or MutableSome constructor function")
	}
}

// MutableNone instantiates an empty MutableOption.
func MutableNone[T any]() OptionalMutable[T] {
	return OptionalMutable[T]{nil, true}
}

// MutableSome instantiates a new OptionP containing the given value.
// The value must exist, i.e. the pointer must not be nil.
func MutableSome[T any](value *T) OptionalMutable[T] {
	if value == nil {
		panic("Cannot create a SomeP out of a nil pointer")
	}
	return OptionalMutable[T]{value, true}
}
