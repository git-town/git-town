package gohacks

import "fmt"

type Option[T fmt.Stringer] struct {
	Value *T
}

func NewEmptyOption[T fmt.Stringer]() Option[T] {
	return Option[T]{nil}
}

func NewOption[T fmt.Stringer](value T) Option[T] {
	return Option[T]{&value}
}

func NewOptionFromPtr[T fmt.Stringer](value *T) Option[T] {
	return Option[T]{value}
}

func (self Option[T]) Get() (T, bool) { //nolint:ireturn
	if self.Has() {
		return *self.Value, true
	}
	var empty T
	return empty, false
}

func (self Option[T]) GetOrDefault() T { //nolint:ireturn
	if self.Has() {
		return *self.Value
	}
	var empty T
	return empty
}

func (self Option[T]) GetOrElse(other T) T { //nolint:ireturn
	if self.Has() {
		return *self.Value
	}
	return other
}

func (self Option[T]) Has() bool {
	return self.Value != nil
}

func (self Option[T]) IsEmpty() bool {
	return !self.Has()
}

func (self Option[T]) MustGet() T { //nolint:ireturn
	if self.Has() {
		return *self.Value
	}
	panic("value not present")
}

func (self Option[T]) String() string {
	return self.StringOr("")
}

func (self Option[T]) StringOr(other string) string {
	if value, has := self.Get(); has {
		return value.String()
	}
	return other
}
