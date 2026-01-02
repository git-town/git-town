package forgedomain

import (
	"strings"

	. "github.com/git-town/git-town/v22/pkg/prelude"
)

// GithubToken is a bearer token to use with the GitHub API.
type GithubToken string

func (self GithubToken) String() string {
	return string(self)
}

func ParseGithubToken(value string) Option[GithubToken] {
	value = strings.TrimSpace(value)
	if value == "" {
		return None[GithubToken]()
	}
	return Some(GithubToken(value))
}
