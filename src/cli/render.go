package cli

// FormatBool converts the given bool into either "yes" or "no".
func FormatBool(value bool) string {
	if value {
		return "yes"
	}
	return "no"
}
