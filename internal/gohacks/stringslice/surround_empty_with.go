package stringslice

// SurroundEmptyWith surrounds all empty strings in the given list with the given character.
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
