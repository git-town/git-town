package forgedomain

import (
	"strings"

	. "github.com/git-town/git-town/v21/pkg/prelude"
)

// GitHubToken is a script that Git Town can call (if configured to do so via GitHubTokenType)
// to retrieve the GitHubToken to use.
type GitHubTokenScript string

func (self GitHubTokenScript) String() string {
	return string(self)
}

func ParseGitHubTokenScript(value string) Option[GitHubTokenScript] {
	value = strings.TrimSpace(value)
	if value == "" {
		return None[GitHubTokenScript]()
	}
	return Some(GitHubTokenScript(value))
}
