package configdomain

import (
	"strings"

	. "github.com/git-town/git-town/v18/pkg/prelude"
)

// GitHubToken is a bearer token to use with the GitHub API.
type CodebergToken string

func (self CodebergToken) String() string {
	return string(self)
}

func ParseCodebergToken(value string) Option[CodebergToken] {
	value = strings.TrimSpace(value)
	if value == "" {
		return None[CodebergToken]()
	}
	return Some(CodebergToken(value))
}
