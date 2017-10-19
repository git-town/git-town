package util

import (
	"strconv"

	"github.com/Originate/exit"
)

// StringToBool parses the given string into a bool
func StringToBool(arg string) bool {
	value, err := strconv.ParseBool(arg)
	exit.If(err)
	return value
}
