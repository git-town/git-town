package stringslice

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
