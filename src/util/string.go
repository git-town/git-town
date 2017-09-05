package util

import (
	"strconv"

	"github.com/Originate/git-town/src/exit"
)

// StringToBool parses the given string into a bool
func StringToBool(arg string) bool {
	value, err := strconv.ParseBool(arg)
	exit.On(err)
	return value
}
