package helpers

// LongestStringLength provides the length of the longest string in the given string collection.
func LongestStringLength(strings []string) (result int) {
	for i := range strings {
		if currentLen := len(strings[i]); currentLen > result {
			result = currentLen
		}
	}
	return result
}
