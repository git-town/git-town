package helpers

// LongestString provides the length of the longest string in the given string collection.
func LongestString(strings []string) (result int) {
	for i := range strings {
		currentLen := len(strings[i])
		if currentLen > result {
			result = currentLen
		}
	}
	return result
}
