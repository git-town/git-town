package cli

import (
	"regexp"
	"strings"
	"sync"
)

var (
	indentOnce sync.Once
	identRE    *regexp.Regexp
)

// Indent outputs the given string with the given level of indentation
// on each line. Each level of indentation is two spaces.
func Indent(message string) string {
	return "  " + strings.ReplaceAll(message, "\n", "\n  ")
}
