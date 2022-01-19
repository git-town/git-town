package cli

import "strings"

// Indent outputs the given string with the given level of indentation
// on each line. Each level of indentation is two spaces.
func Indent(message string) string {
	return "  " + strings.ReplaceAll(message, "\n", "\n  ")
}
