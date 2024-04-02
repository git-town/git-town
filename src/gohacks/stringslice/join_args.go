package stringslice

import (
	"fmt"
	"strings"
)

// JoinArgs provides the given command arguments joined into a single string
// while quote-escaping as needed.
func JoinArgs(args []string) string {
	quoted := make([]string, len(args))
	for a, arg := range args {
		if arg == "" {
			quoted[a] = `""`
		} else if strings.ContainsRune(arg, ' ') {
			quoted[a] = fmt.Sprintf("%q", arg)
		} else {
			quoted[a] = arg
		}
	}
	return strings.Join(quoted, " ")
}
