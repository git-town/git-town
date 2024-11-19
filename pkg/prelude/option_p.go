package prelude

import (
	"encoding/json"
	"fmt"
)

// OptionalMut represents a value that is both optional and mutable.
type OptionalMut[T any] struct {
	Value *T
}

// Get provides null-safe mutable access to the contained value.
func (self OptionalMut[T]) Get() (value *T, hasValue bool) {
	if self.IsSome() {
		return self.Value, true
	}
	return nil, false
}

// GetOrPanic provides unsafe mutable access to the contained value.
func (self OptionalMut[T]) GetOrPanic() *T {
	if value, has := self.Get(); has {
		return value
	}
	panic("value not present")
}

// IsNone indicates whether this option instance contains nothing.
func (self OptionalMut[T]) IsNone() bool {
	return self.Value == nil
}

// IsSome indicates whether this option instance contains a value.
func (self OptionalMut[T]) IsSome() bool {
	return self.Value != nil
}

// MarshalJSON is used when serializing this type to JSON.
func (self OptionalMut[T]) MarshalJSON() ([]byte, error) {
	if value, hasValue := self.Get(); hasValue {
		return json.Marshal(*value)
	}
	return json.Marshal(nil)
}

// String provides the string serialization of the contained value.
// None gets serialized into an empty string.
func (self OptionalMut[T]) String() string {
	return self.StringOr("")
}

// StringOr provideds the string serialization of the contained value.
// None gets serialized into the given alternative string representation.
func (self OptionalMut[T]) StringOr(other string) string {
	if self.IsSome() {
		return fmt.Sprint(self.Value)
	}
	return other
}

// ToOption provides an immutable copy of this OptionalMut.
func (self OptionalMut[T]) ToOption() Option[T] {
	if value, hasValue := self.Get(); hasValue {
		return Some(*value)
	}
	return None[T]()
}

// UnmarshalJSON is used when de-serializing this type from JSON.
func (self *OptionalMut[T]) UnmarshalJSON(b []byte) error {
	if string(b) == "null" {
		self.Value = nil
		return nil
	}
	var value T
	self.Value = &value
	return json.Unmarshal(b, &self.Value)
}
