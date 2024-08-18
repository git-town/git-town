package configdomain

import (
	"strings"

	. "github.com/git-town/git-town/v15/pkg/prelude"
)

type GitUserName string

func (self GitUserName) String() string {
	return string(self)
}

func ParseGitUserName(value string) Option[GitUserName] {
	value = strings.TrimSpace(value)
	if value == "" {
		return None[GitUserName]()
	}
	return Some(GitUserName(value))
}
