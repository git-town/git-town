package prelude

import "fmt"

// None instantiates an Option of the given type containing nothing.
func None[T fmt.Stringer]() Option[T] {
	return Option[T]{nil}
}
