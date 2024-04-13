package execute

import (
	"fmt"

	"github.com/git-town/git-town/v14/src/git/gitdomain"
)

// FailureCollector helps avoid excessive error checking
// while gathering a larger number of values through fallible operations.
// This is based on ideas outlined in https://go.dev/blog/errors-are-values.
//
// Please be aware that using this technique can lead to executing logic that would normally not run,
// using potentially invalid data, and potentially leading to unexpected runtime exceptions and side effects.
// Use with care and only if it's abundantly clear and obvious that there are no negative side effects.
// This is an anti-pattern in code to work arount an anti-pattern in the language.
type FailureCollector struct {
	Err error `exhaustruct:"optional"`
}

func (self *FailureCollector) BranchInfos(value gitdomain.BranchInfos, err error) gitdomain.BranchInfos {
	self.Check(err)
	return value
}

// Check registers the given error and indicates
// whether this ErrorChecker contains an error now.
func (self *FailureCollector) Check(err error) bool {
	if self.Err == nil {
		self.Err = err
	}
	return self.Err != nil
}

// Fail registers the error constructed using the given format arguments.
func (self *FailureCollector) Fail(format string, a ...any) {
	self.Check(fmt.Errorf(format, a...))
}

// Remotes provides the domain.Remotes part of the given fallible function result
// while registering the given error.
func (self *FailureCollector) Remotes(value gitdomain.Remotes, err error) gitdomain.Remotes {
	self.Check(err)
	return value
}
