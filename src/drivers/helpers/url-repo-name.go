package helpers

import (
	"regexp"
	"strings"
)

// GetURLRepositoryName returns the repository name contains within the given Git URL.
func GetURLRepositoryName(url string) string {
	hostname := URLHostname(url)
	repositoryNameRegex := regexp.MustCompile(".*" + hostname + "[/:](.+)")
	matches := repositoryNameRegex.FindStringSubmatch(url)
	if matches == nil {
		return ""
	}
	return strings.TrimSuffix(matches[1], ".git")
}
