package helpers

import (
	"regexp"
	"strings"
)

// URLRepositoryName provides the repository name contains within the given Git URL.
func URLRepositoryName(url string) string {
	hostname := URLHostname(url)
	repositoryNameRegex := regexp.MustCompile(".*" + hostname + "[/:](.+)")
	matches := repositoryNameRegex.FindStringSubmatch(url)
	if matches == nil {
		return ""
	}
	return strings.TrimSuffix(matches[1], ".git")
}
