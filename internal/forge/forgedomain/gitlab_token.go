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

func ParseGitlabToken(value string) Option[GitlabToken] {
	value = strings.TrimSpace(value)
	if value == "" {
		return None[GitlabToken]()
	}
	return Some(GitlabToken(value))
}
