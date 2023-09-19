package persistence

import (
	"regexp"
	"strings"
)

func SanitizePath(dir string) string {
	replaceCharacterRE := regexp.MustCompile("[[:^alnum:]]")
	sanitized := replaceCharacterRE.ReplaceAllString(dir, "-")
	sanitized = strings.ToLower(sanitized)
	replaceDoubleMinusRE := regexp.MustCompile("--+") // two or more dashes
	sanitized = replaceDoubleMinusRE.ReplaceAllString(sanitized, "-")
	for strings.HasPrefix(sanitized, "-") {
		sanitized = sanitized[1:]
	}
	return sanitized
}
