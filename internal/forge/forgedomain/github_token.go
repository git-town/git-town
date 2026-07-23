package forgedomain

import (
	"github.com/git-town/git-town/v24/internal/gohacks/stringss"
	. "github.com/git-town/git-town/v24/pkg/prelude"
)

// GithubToken is a bearer token to use with the GitHub API.
type GithubToken stringss.Trimmed

func (self GithubToken) String() string {
	return string(self)
}

func ParseGithubToken(value stringss.Trimmed) Option[GithubToken] {
	if value == "" {
		return None[GithubToken]()
	}
	return Some(GithubToken(value))
}
