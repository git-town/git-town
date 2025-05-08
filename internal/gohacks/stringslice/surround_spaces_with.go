package stringslice

import "strings"

// SurroundSpacesWith surrounds all strings that contain a space in the given list with the given character.
func SurroundSpacesWith(args []string, surround string) []string {
	result := make([]string, len(args))
	for t, text := range args {
		if strings.Contains(text, " ") {
			result[t] = surround + text + surround
		} else {
			result[t] = text
		}
	}
	return result
}
