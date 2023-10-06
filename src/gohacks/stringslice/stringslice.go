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
	if text == "" {
		return []string{}
	}
	return strings.Split(text, "\n")
}

// Longest provides the length of the longest string in the given string slice.
func Longest(strings []string) int {
	result := 0
	for _, s := range strings {
		if currentLen := len(s); currentLen > result {
			result = currentLen
		}
	}
	return result
}

// SurroundEmptyWith surrounds all empty strings in the given list with the given character
func SurroundEmptyWith(strings []string, surround string) []string {
	result := make([]string, len(strings))
	for t, text := range strings {
		if text == "" {
			result[t] = surround + surround
		} else {
			result[t] = text
		}
	}
	return result
}
