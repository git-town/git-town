package forgedomain

import (
	"strings"

	. "github.com/git-town/git-town/v22/pkg/prelude"
)

// GitlabToken is a bearer token to use with the GitLab API.
type GitlabToken string

func (self GitlabToken) String() string {
	return string(self)
}

func ParseGitlabToken(valueOpt Option[string]) Option[GitlabToken] {
	if value, has := valueOpt.Get(); has {
		return Some(GitlabToken(strings.TrimSpace(value)))
	}
	return None[GitlabToken]()
}
