package util

import (
	"fmt"
	"strconv"
)

// StringToBool parses the given string into a bool.
func StringToBool(arg string) bool {
	value, err := strconv.ParseBool(arg)
	if err != nil {
		fmt.Println("Cannot ")
	}
	return value
}
