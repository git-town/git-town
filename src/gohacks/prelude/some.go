package prelude

import "fmt"

// Some instantiates a new option containing the given value.
func Some[T fmt.Stringer](value T) Option[T] {
	return Option[T]{&value}
}
