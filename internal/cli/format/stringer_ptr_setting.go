package format

import (
	"fmt"

	. "github.com/git-town/git-town/v23/pkg/prelude"
)

// OptionalStringerSetting provides a printable version of the given configuration value.
// The configuration value must conform to the fmt.Stringer interface.
func OptionalStringerSetting[T fmt.Stringer](option Option[T]) string {
	return option.StringOr("(not set)")
}

// ConfiguredStringerSetting provides a printable version of the given configuration value without exposing its contents.
func ConfiguredStringerSetting[T fmt.Stringer](option Option[T]) string {
	if option.IsSome() {
		return "(configured)"
	}
	return OptionalStringerSetting(option)
}
