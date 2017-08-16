package util

import (
	"strconv"

	"github.com/Originate/git-town/src/logs"
)

// StringToBool parses the given string into a bool
func StringToBool(arg string) bool {
	value, err := strconv.ParseBool(arg)
	logs.FatalOn(err)
	return value
}
