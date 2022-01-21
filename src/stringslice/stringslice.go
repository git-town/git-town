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
