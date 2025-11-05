package filesystem

import (
	"strings"
	"unicode"
)

// FolderName converts the given scenario name into a string
// that can be used safely as a folder name on the filesystem.
func FolderName(scenarioName string) string {
	result := strings.Builder{}
	lastRune := ' '
	for _, r := range scenarioName {
		if unicode.IsLetter(r) {
			r = unicode.ToLower(r)
			result.WriteRune(r)
			lastRune = r
			continue
		}
		if lastRune != '_' {
			result.WriteRune('_')
			lastRune = '_'
		}
	}
	return result.String()
}
