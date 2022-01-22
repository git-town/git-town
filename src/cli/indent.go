package cli

import (
	"regexp"
	"strings"
	"sync"
)

var indentOnce sync.Once
var identRE *regexp.Regexp

// Indent outputs the given string with the given level of indentation
// on each line. Each level of indentation is two spaces.
func Indent(message string) string {
	result := "  " + strings.ReplaceAll(message, "\n", "\n  ")
	indentOnce.Do(func() { identRE = regexp.MustCompile("\n  \n") })
	return identRE.ReplaceAllString(result, "\n\n")
}
