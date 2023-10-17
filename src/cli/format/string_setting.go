package format

func StringSetting(text string) string {
	if text == "" {
		return "(not set)"
	}
	return text
}
