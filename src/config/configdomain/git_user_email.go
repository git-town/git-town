package configdomain

import (
	"strings"

	. "github.com/git-town/git-town/v14/src/gohacks/prelude"
)

type GitUserEmail string

func (self GitUserEmail) String() string {
	return string(self)
}

func NewGitUserEmailOption(value string) Option[GitUserEmail] {
	value = strings.TrimSpace(value)
	if value == "" {
		return None[GitUserEmail]()
	}
	return Some(GitUserEmail(value))
}
