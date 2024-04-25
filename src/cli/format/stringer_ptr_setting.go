package format

import (
	"fmt"

	. "github.com/git-town/git-town/v14/src/gohacks/prelude"
)

// OptionalStringerSetting provides a printable version of the given configuration value.
// The configuration value must conform to the fmt.Stringer interface.
func OptionalStringerSetting[T fmt.Stringer](option Option[T]) string {
	return option.StringOr("(not set)")
}
