package configdomain

import (
	"strings"

	. "github.com/git-town/git-town/v14/src/gohacks/prelude"
)

// GitLabToken is a bearer token to use with the GitLab API.
type GitLabToken string

func (self GitLabToken) String() string {
	return string(self)
}

func NewGitLabToken(value string) GitLabToken {
	value = strings.TrimSpace(value)
	if value == "" {
		panic("empty GitHub token")
	}
	return GitLabToken(value)
}

func NewGitLabTokenOption(value string) Option[GitLabToken] {
	value = strings.TrimSpace(value)
	if value == "" {
		return None[GitLabToken]()
	}
	return Some(NewGitLabToken(value))
}
