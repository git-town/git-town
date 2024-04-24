package format

import "fmt"

// StringSetting provides a printable version of the given string configuration value.
func StringerPtrSetting(token fmt.Stringer) string {
	if token == nil {
		return "(not set)"
	}
	return token.String()
}
