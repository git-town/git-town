package execute

import (
	"github.com/git-town/git-town/v14/src/git/gitdomain"
	"github.com/git-town/git-town/v14/src/gohacks"
)

// FailureCollector is an ErrorCollector wrapper that provides convenience functions to check fallible domain types.
//
// Please be aware that using this technique can lead to executing logic that would normally not run,
// using potentially invalid data, and potentially leading to unexpected runtime exceptions and side effects.
// Use with care and only if it's abundantly clear and obvious that there are no negative side effects.
// This is an anti-pattern in code to work arount an anti-pattern in the language.
type FailureCollector struct {
	gohacks.ErrorCollector `exhaustruct:"optional"`
}

func (self *FailureCollector) BranchInfos(value gitdomain.BranchInfos, err error) gitdomain.BranchInfos {
	self.Check(err)
	return value
}

// Remotes provides the domain.Remotes part of the given fallible function result
// while registering the given error.
func (self *FailureCollector) Remotes(value gitdomain.Remotes, err error) gitdomain.Remotes {
	self.Check(err)
	return value
}
