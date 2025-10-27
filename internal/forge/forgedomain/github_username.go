package forgedomain

import (
	"strings"

	. "github.com/git-town/git-town/v22/pkg/prelude"
)

// GitHubUsername is a GitHub username to use as a prefix for branch names.
type GitHubUsername string

func (self GitHubUsername) String() string {
	return string(self)
}

func ParseGitHubUsername(value string) Option[GitHubUsername] {
	value = strings.TrimSpace(value)
	if value == "" {
		return None[GitHubUsername]()
	}
	return Some(GitHubUsername(value))
}
