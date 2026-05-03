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

func ParseGithubToken(valueOpt Option[string]) Option[GithubToken] {
	if value, has := valueOpt.Get(); has {
		return Some(GithubToken(strings.TrimSpace(value)))
	}
	return None[GithubToken]()
}
