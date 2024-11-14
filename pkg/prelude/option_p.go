package prelude

import (
	"encoding/json"
	"fmt"
)

// MutableOption ("option for a pointer value") is an Option
// that provides direct access to the encapsulated value
// by storing and providing a pointer to the value.
// This is useful for mutable or singleton values,
// or values that are too large to copy around all the time.
type MutableOption[T any] struct {
	Value *T
}

// Get provides a copy of the contained value
// as well as an indicator whether that value exists.
func (self MutableOption[T]) Get() (value *T, hasValue bool) {
	if self.IsSome() {
		return self.Value, true
	}
	return nil, false
}

// GetOrPanic provides a copy of the contained value.
// Panics if this option contains nothing.
func (self MutableOption[T]) GetOrPanic() *T {
	if value, has := self.Get(); has {
		return value
	}
	panic("value not present")
}

// IsNone indicates whether this option instance contains nothing.
func (self MutableOption[T]) IsNone() bool {
	return self.Value == nil
}

// IsSome indicates whether this option instance contains a value.
func (self MutableOption[T]) IsSome() bool {
	return self.Value != nil
}

// MarshalJSON is used when serializing this OptionP to JSON.
func (self MutableOption[T]) MarshalJSON() ([]byte, error) {
	if value, hasValue := self.Get(); hasValue {
		return json.Marshal(*value)
	}
	return json.Marshal(nil)
}

// String provides the string serialization of the contained value.
// If this option contains nothing, you get an empty string.
func (self MutableOption[T]) String() string {
	return self.StringOr("")
}

// StringOr provideds the string serialization of the contained value.
// If this option contains nothing, you get the given alternative string representation.
func (self MutableOption[T]) StringOr(other string) string {
	if self.IsSome() {
		return fmt.Sprint(self.Value)
	}
	return other
}

// converts this OptionP to an Option
func (self MutableOption[T]) ToOption() Option[T] {
	if value, hasValue := self.Get(); hasValue {
		return Some(*value)
	}
	return None[T]()
}

// UnmarshalJSON is used when de-serializing JSON into an OptionP.
func (self *MutableOption[T]) UnmarshalJSON(b []byte) error {
	if string(b) == "null" {
		self.Value = nil
		return nil
	}
	var value T
	self.Value = &value
	return json.Unmarshal(b, &self.Value)
}
