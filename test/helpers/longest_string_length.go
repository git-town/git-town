package helpers

// LongestStringLength provides the length of the longest string in the given string collection.
func LongestStringLength(strings []string) (result int) {
	for _, s := range strings {
		if currentLen := len(s); currentLen > result {
			result = currentLen
		}
	}
	return result
}
