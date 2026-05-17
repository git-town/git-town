package forgedomain

import (
	"github.com/git-town/git-town/v23/internal/gohacks/stringss"
	. "github.com/git-town/git-town/v23/pkg/prelude"
)

// GithubToken is a bearer token to use with the GitHub API.
type GithubToken stringss.TrimmedString

func (self GithubToken) String() string {
	return string(self)
}

func ParseGithubToken(value stringss.TrimmedString) Option[GithubToken] {
	if value == "" {
		return None[GithubToken]()
	}
	return Some(GithubToken(value))
}
