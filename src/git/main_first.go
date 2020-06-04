package git

// MainFirst provides the given list of strings, with an element "main" moved to its first position.
func MainFirst(input []string) []string {
	result := make([]string, 0, len(input))
	hasMain := false
	for i := range input {
		if input[i] == "main" {
			hasMain = true
		} else {
			result = append(result, input[i])
		}
	}
	if hasMain {
		result = append([]string{"main"}, result...)
	}
	return result
}
