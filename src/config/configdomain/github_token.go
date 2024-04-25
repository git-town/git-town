package configdomain

import (
	"strings"

	. "github.com/git-town/git-town/v14/src/gohacks/prelude"
)

// GitHubToken is a bearer token to use with the GitHub API.
type GitHubToken string

func (self GitHubToken) String() string {
	return string(self)
}

func NewGitHubToken(value string) GitHubToken {
	value = strings.TrimSpace(value)
	if value == "" {
		panic("empty GitHub token")
	}
	return GitHubToken(value)
}

func NewGitHubTokenOption(value string) Option[GitHubToken] {
	value = strings.TrimSpace(value)
	if value == "" {
		return None[GitHubToken]()
	}
	return Some(NewGitHubToken(value))
}
