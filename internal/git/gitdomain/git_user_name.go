package gitdomain

import (
	"strings"

	. "github.com/git-town/git-town/v23/pkg/prelude"
)

type GitUserName string

func (self GitUserName) String() string {
	return string(self)
}

func GitUserNameFromString(value string) Option[GitUserName] {
	value = strings.TrimSpace(value)
	if value == "" {
		return None[GitUserName]()
	}
	return Some(GitUserName(value))
}

func GitUserNameFromStringOpt(valueOpt Option[string]) Option[GitUserName] {
	if value, has := valueOpt.Get(); has {
		return Some(GitUserName(value))
	}
	return None[GitUserName]()
}
