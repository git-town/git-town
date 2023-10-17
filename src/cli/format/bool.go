package format

// Bool converts the given bool into either "yes" or "no".
func Bool(value bool) string {
	if value {
		return "yes"
	}
	return "no"
}
