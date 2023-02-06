// Package stringslice provides helper functions for working with slices of strings.
package stringslice

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

// Hoist provides the given list of strings, with the given element moved to the first position.
func Hoist(list []string, element string) []string {
	result := make([]string, 0, len(list))
	hasMain := false
	for _, input := range list {
		if input == element {
			hasMain = true
		} else {
			result = append(result, input)
		}
	}
	if hasMain {
		result = append([]string{"main"}, result...)
	}
	return result
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
