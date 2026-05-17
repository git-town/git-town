package forgedomain

import (
	"github.com/git-town/git-town/v23/internal/gohacks/stringss"
	. "github.com/git-town/git-town/v23/pkg/prelude"
)

// GitlabToken is a bearer token to use with the GitLab API.
type GitlabToken string

func (self GitlabToken) String() string {
	return string(self)
}

func ParseGitlabToken(value stringss.TrimmedString) Option[GitlabToken] {
	if value == "" {
		return None[GitlabToken]()
	}
	return Some(GitlabToken(value))
}
