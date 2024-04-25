package format

import (
	"fmt"

	"github.com/git-town/git-town/v14/src/gohacks"
)

// OptionalStringerSetting provides a printable version of the given configuration value.
// The configuration value must conform to the fmt.Stringer interface.
func OptionalStringerSetting[T fmt.Stringer](option gohacks.Option[T]) string {
	return option.StringOr("(not set)")
}
