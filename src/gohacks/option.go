package gohacks

import "fmt"

type Option[T fmt.Stringer] struct {
	Value *T
}

func NewOption[T fmt.Stringer](value T) Option[T] {
	return Option[T]{&value}
}

func NewOptionNone[T fmt.Stringer]() Option[T] {
	return Option[T]{nil}
}

func (self Option[T]) Get() (T, bool) { //nolint:ireturn
	if self.IsSome() {
		return *self.Value, true
	}
	var empty T
	return empty, false
}

func (self Option[T]) GetOrDefault() T { //nolint:ireturn
	if value, has := self.Get(); has {
		return value
	}
	var empty T
	return empty
}

func (self Option[T]) GetOrElse(other T) T { //nolint:ireturn
	if value, has := self.Get(); has {
		return value
	}
	return other
}

func (self Option[T]) GetOrPanic() T { //nolint:ireturn
	if value, has := self.Get(); has {
		return value
	}
	panic("value not present")
}

func (self Option[T]) IsNone() bool {
	return self.Value == nil
}

func (self Option[T]) IsSome() bool {
	return self.Value != nil
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
