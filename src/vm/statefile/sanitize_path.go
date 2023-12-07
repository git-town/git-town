package statefile

import (
	"regexp"
	"strings"

	"github.com/git-town/git-town/v11/src/domain"
)

func SanitizePath(dir domain.RepoRootDir) string {
	replaceCharacterRE := regexp.MustCompile("[[:^alnum:]]")
	sanitized := replaceCharacterRE.ReplaceAllString(dir.String(), "-")
	sanitized = strings.ToLower(sanitized)
	replaceDoubleMinusRE := regexp.MustCompile("--+") // two or more dashes
	sanitized = replaceDoubleMinusRE.ReplaceAllString(sanitized, "-")
	for strings.HasPrefix(sanitized, "-") {
		sanitized = sanitized[1:]
	}
	return sanitized
}
