package gitdomain

import (
	"strings"

	. "github.com/git-town/git-town/v21/pkg/prelude"
)

type GitUserEmail string

func (self GitUserEmail) String() string {
	return string(self)
}

func ParseGitUserEmail(value string) Option[GitUserEmail] {
	value = strings.TrimSpace(value)
	if value == "" {
		return None[GitUserEmail]()
	}
	return Some(GitUserEmail(value))
}
