package gitdomain

import (
	"github.com/git-town/git-town/v23/internal/gohacks/stringss"
	. "github.com/git-town/git-town/v23/pkg/prelude"
)

type GitUserEmail stringss.TrimmedString

func (self GitUserEmail) String() string {
	return string(self)
}

func ParseGitUserEmail(value stringss.TrimmedString) Option[GitUserEmail] {
	if value == "" {
		return None[GitUserEmail]()
	}
	return Some(GitUserEmail(value))
}

func ParseGitUserEmailOpt(valueOpt Option[string]) Option[GitUserEmail] {
	if value, has := valueOpt.Get(); has {
		return Some(GitUserEmail(value))
	}
	return None[GitUserEmail]()
}
