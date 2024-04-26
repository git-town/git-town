package configdomain

import (
	"strings"

	. "github.com/git-town/git-town/v14/src/gohacks/prelude"
)

type GitUserEmail string

func (self GitUserEmail) String() string {
	return string(self)
}

func NewGitUserEmail(value string) GitUserEmail {
	value = strings.TrimSpace(value)
	if value == "" {
		panic("empty Git user email")
	}
	return GitUserEmail(value)
}

func NewGitUserEmailOption(value string) Option[GitUserEmail] {
	value = strings.TrimSpace(value)
	if value == "" {
		return None[GitUserEmail]()
	}
	return Some(NewGitUserEmail(value))
}
