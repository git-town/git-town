package forgedomain

import (
	"github.com/git-town/git-town/v24/internal/gohacks/stringss"
	. "github.com/git-town/git-town/v24/pkg/prelude"
)

// GitlabToken is a bearer token to use with the GitLab API.
type GitlabToken stringss.Trimmed

func (self GitlabToken) String() string {
	return string(self)
}

func ParseGitlabToken(value stringss.Trimmed) Option[GitlabToken] {
	if value == "" {
		return None[GitlabToken]()
	}
	return None[GitlabToken]()
}
