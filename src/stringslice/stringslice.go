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

// Last provides a pointer to the last element in the slice.
func Last(list []string) *string {
	idx := len(list)
	if idx == 0 {
		return nil
	}
	return &list[idx-1]
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

func RemoveMany(list []string, toRemoves []string) []string {
	result := []string{}
	for _, element := range list {
		if !Contains(toRemoves, element) {
			result = append(result, element)
		}
	}
	return result
}
