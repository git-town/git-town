package util

import (
	"fmt"
	"os"
	"strconv"
)

// StringToBool parses the given string into a bool
func StringToBool(arg string) bool {
	value, err := strconv.ParseBool(arg)
	if err != nil {
		fmt.Printf("cannot convert string %q to bool: %v\n", arg, err)
		os.Exit(1)
	}
	return value
}
