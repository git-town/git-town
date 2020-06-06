package util

import (
	"strings"
)

// DoesStringArrayContain returns whether the given string slice
// contains the given string.
func DoesStringArrayContain(list []string, value string) bool {
	for _, element := range list {
		if element == value {
			return true
		}
	}
	return false
}

// Indent outputs the given string with the given level of indentation
// on each line. Each level of indentation is two spaces.
func Indent(message string) string {
	return "  " + strings.Replace(message, "\n", "\n  ", -1)
}

// RemoveStringFromSlice returns a new string slice which is the given string slice
// with the given string removed.
func RemoveStringFromSlice(list []string, value string) (result []string) {
	for _, element := range list {
		if element != value {
			result = append(result, element)
		}
	}
	return
}
