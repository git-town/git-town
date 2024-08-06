package configdomain

import (
	"strings"

	. "github.com/git-town/git-town/v15/internal/gohacks/prelude"
)

// GitHubToken is a bearer token to use with the GitHub API.
type GitHubToken string

func (self GitHubToken) String() string {
	return string(self)
}

func ParseGitHubToken(value string) Option[GitHubToken] {
	value = strings.TrimSpace(value)
	if value == "" {
		return None[GitHubToken]()
	}
	return Some(GitHubToken(value))
}
