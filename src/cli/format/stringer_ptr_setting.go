package format

import (
	"fmt"

	"github.com/git-town/git-town/v14/src/gohacks"
)

// StringSetting provides a printable version of the given string configuration value.
func StringerPtrSetting(token fmt.Stringer) string {
	if gohacks.IsNil(token) {
		return "(not set)"
	}
	return token.String()
}
