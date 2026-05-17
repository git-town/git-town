package gitdomain

import (
	"github.com/git-town/git-town/v23/internal/gohacks/stringss"
	. "github.com/git-town/git-town/v23/pkg/prelude"
)

type GitUserName stringss.Trimmed

func (self GitUserName) String() string {
	return string(self)
}

func GitUserNameOptFromString(value stringss.Trimmed) Option[GitUserName] {
	if value == "" {
		return None[GitUserName]()
	}
	return Some(GitUserName(value))
}

func GitUserNameOptFromStringOpt(valueOpt Option[string]) Option[GitUserName] {
	if value, has := valueOpt.Get(); has {
		return Some(GitUserName(value))
	}
	return None[GitUserName]()
}
