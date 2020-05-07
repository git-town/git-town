package util

import (
	"log"
	"strconv"
)

// StringToBool parses the given string into a bool
func StringToBool(arg string) bool {
	value, err := strconv.ParseBool(arg)
	if err != nil {
		log.Fatal(err)
	}
	return value
}
