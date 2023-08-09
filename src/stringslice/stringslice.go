// Package stringslice provides helper functions for working with slices of strings.
package stringslice

import (
	"fmt"
	"strings"
)

func AppendAllMissing(list []string, elements []string) []string {
	for _, element := range elements {
		if !Contains(list, element) {
			list = append(list, element)
		}
	}
	return list
}

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

// Contains returns whether the given string slice
// contains the given string.
func Contains(list []string, value string) bool {
	for _, element := range list {
		if element == value {
			return true
		}
	}
	return false
}

func FirstElementOr(list []string, alternative string) string {
	if len(list) > 0 {
		return list[0]
	}
	return alternative
}

// Hoist provides the given list of strings, with the given element moved to the first position.
func Hoist(list []string, element string) []string {
	result := make([]string, 0, len(list))
	hasElement := false
	for _, input := range list {
		if input == element {
			hasElement = true
		} else {
			result = append(result, input)
		}
	}
	if hasElement {
		result = append([]string{element}, result...)
	}
	return result
}

func Lines(text string) []string {
	return strings.Split(text, "\n")
}

// Remove returns a new string slice which is the given string slice
// with the given string removed.
func Remove(list []string, value string) []string {
	result := make([]string, 0, len(list)-1)
	for _, element := range list {
		if element != value {
			result = append(result, element)
		}
	}
	return result
}
