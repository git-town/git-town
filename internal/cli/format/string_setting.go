package format

// StringSetting provides a printable version of the given string configuration value.
func StringSetting(text string) string {
	if text == "" {
		return "(not set)"
	}
	return text
}
