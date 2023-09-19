package gohacks

import (
	"fmt"

	"github.com/git-town/git-town/v9/src/config"
	"github.com/git-town/git-town/v9/src/domain"
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

// Check registers the given error and indicates
// whether this ErrorChecker contains an error now.
func (ec *FailureCollector) Check(err error) bool {
	if ec.Err == nil {
		ec.Err = err
	}
	return ec.Err != nil
}

// Fail registers the error constructed using the given format arguments.
func (ec *FailureCollector) Fail(format string, a ...any) {
	ec.Check(fmt.Errorf(format, a...))
}

// Bool provides the bool part of the given fallible function result
// while registering the given error.
func (ec *FailureCollector) Bool(value bool, err error) bool {
	ec.Check(err)
	return value
}

func (ec *FailureCollector) Branches(value domain.Branches, err error) domain.Branches {
	ec.Check(err)
	return value
}

func (ec *FailureCollector) BranchesSyncStatus(value domain.BranchInfos, err error) domain.BranchInfos {
	ec.Check(err)
	return value
}

// Hosting provides the config.Hosting part of the given fallible function result
// while registering the given error.
func (ec *FailureCollector) Hosting(value config.Hosting, err error) config.Hosting {
	ec.Check(err)
	return value
}

// PullBranchStrategy provides the string part of the given fallible function result
// while registering the given error.
func (ec *FailureCollector) PullBranchStrategy(value config.PullBranchStrategy, err error) config.PullBranchStrategy {
	ec.Check(err)
	return value
}

// Remotes provides the domain.Remotes part of the given fallible function result
// while registering the given error.
func (ec *FailureCollector) Remotes(value domain.Remotes, err error) domain.Remotes {
	ec.Check(err)
	return value
}

// String provides the string part of the given fallible function result
// while registering the given error.
func (ec *FailureCollector) String(value string, err error) string {
	ec.Check(err)
	return value
}

// Strings provides the []string part of the given fallible function result
// while registering the given error.
func (ec *FailureCollector) Strings(value []string, err error) []string {
	ec.Check(err)
	return value
}

// SyncStrategy provides the string part of the given fallible function result
// while registering the given error.
func (ec *FailureCollector) SyncStrategy(value config.SyncStrategy, err error) config.SyncStrategy {
	ec.Check(err)
	return value
}
