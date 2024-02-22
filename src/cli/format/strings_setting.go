package format

// StringSetting provides a printable version of the given string configuration value.
func StringsSetting(text string) string {
	if text == "" {
		return "(none)"
	}
	return text
}
