package helpers

import "unicode"

// FolderName converts the given scenario name into a string
// that can be used safely as a folder name on the filesystem.
func FolderName(scenarioName string) string {
	result := ""
	lastRune := ' '
	for _, r := range scenarioName {
		if unicode.IsLetter(r) {
			r = unicode.ToLower(r)
			result += string(r)
			lastRune = r
			continue
		}
		if lastRune != '_' {
			result += "_"
			lastRune = '_'
		}
	}
	return result
}
