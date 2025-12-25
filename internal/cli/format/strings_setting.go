package format

import "github.com/git-town/git-town/v22/internal/messages"

// StringsSetting provides a printable version of the given []string configuration value.
func StringsSetting(text string) string {
	if text == "" {
		return messages.DialogResultNone
	}
	return text
}
