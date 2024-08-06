package stringslice

import (
	"strings"
)

// JoinArgs provides the given command arguments joined into a single string
// while quote-escaping as needed.
func JoinArgs(args []string) string {
	quoted := make([]string, len(args))
	for a, arg := range args {
		switch {
		case arg == "":
			quoted[a] = `""`
		case strings.ContainsRune(arg, ' '):
			quoted[a] = `"` + arg + `"`
		default:
			quoted[a] = arg
		}
	}
	return strings.Join(quoted, " ")
}
