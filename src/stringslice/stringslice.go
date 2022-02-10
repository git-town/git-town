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

// MainFirst provides the given list of strings, with an element "main" moved to its first position.
func MainFirst(inputs []string) []string {
	result := make([]string, 0, len(inputs))
	hasMain := false
	for _, input := range inputs {
		if input == "main" {
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
func Remove(list []string, value string) (result []string) {
	for _, element := range list {
		if element != value {
			result = append(result, element)
		}
	}
	return
}
