package prelude

import "fmt"

// NewOptionNone instantiates a new option containing nothing.
func None[T fmt.Stringer]() Option[T] {
	return Option[T]{nil}
}
