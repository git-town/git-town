// Package stringslice provides helper functions for working with slices of strings.
package stringslice

import (
	"fmt"
	"strings"
)

// Connect provides a human-readable serialization of the given strings list.
func Connect(list []string) string {
	count := len(list)
	if count == 0 {
		return ""
	}
	if count == 1 {
		return quote(list[0])
	}
	if count == 2 {
		return fmt.Sprintf("%q and %q", list[0], list[1])
	}
	result := quote(list[0])
	for i, element := range list {
		if i == 0 || i == count-1 {
			continue
		}
		result = result + ", " + quote(element)
	}
	return result + ", and " + quote(list[count-1])
}

func quote(text string) string {
	return "\"" + text + "\""
}

func Lines(text string) []string {
	return strings.Split(text, "\n")
}
