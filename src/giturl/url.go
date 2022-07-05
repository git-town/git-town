// Package giturl provides facilities to work with the special URL formats used in Git remotes.
package giturl

import (
	"regexp"
	"strings"
)

// Host provides the hostname contained within the given Git hosting URL.
func Host(url string) string {
	hostnameRegex := regexp.MustCompile("(^[^:]*://([^@]*@)?|[^@]*@)([^/:]+).*")
	matches := hostnameRegex.FindStringSubmatch(url)
	if matches == nil {
		return ""
	}
	return matches[3]
}

// Repo provides the repository name contained within the given Git hosting URL.
func Repo(url string) string {
	hostname := Host(url)
	repositoryNameRegex := regexp.MustCompile(".*" + hostname + "[/:](.+)")
	matches := repositoryNameRegex.FindStringSubmatch(url)
	if matches == nil {
		return ""
	}
	return strings.TrimSuffix(matches[1], ".git")
}
