package gohacks

import (
	"fmt"

	"github.com/git-town/git-town/v9/src/config"
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/git"
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
func (fc *FailureCollector) Check(err error) bool {
	if fc.Err == nil {
		fc.Err = err
	}
	return fc.Err != nil
}

// Fail registers the error constructed using the given format arguments.
func (fc *FailureCollector) Fail(format string, a ...any) {
	fc.Check(fmt.Errorf(format, a...))
}

// Bool provides the bool part of the given fallible function result
// while registering the given error.
func (fc *FailureCollector) Bool(value bool, err error) bool {
	fc.Check(err)
	return value
}

func (fc *FailureCollector) Branches(value domain.Branches, err error) domain.Branches {
	fc.Check(err)
	return value
}

func (fc *FailureCollector) BranchesSyncStatus(value domain.BranchInfos, err error) domain.BranchInfos {
	fc.Check(err)
	return value
}

// Hosting provides the config.Hosting part of the given fallible function result
// while registering the given error.
func (fc *FailureCollector) Hosting(value config.Hosting, err error) config.Hosting {
	fc.Check(err)
	return value
}

// PullBranchStrategy provides the string part of the given fallible function result
// while registering the given error.
func (fc *FailureCollector) PullBranchStrategy(value config.PullBranchStrategy, err error) config.PullBranchStrategy {
	fc.Check(err)
	return value
}

// Remotes provides the domain.Remotes part of the given fallible function result
// while registering the given error.
func (fc *FailureCollector) Remotes(value domain.Remotes, err error) domain.Remotes {
	fc.Check(err)
	return value
}

func (fc *FailureCollector) RepoStatus(value git.RepoStatus, err error) git.RepoStatus {
	fc.Check(err)
	return value
}

// String provides the string part of the given fallible function result
// while registering the given error.
func (fc *FailureCollector) String(value string, err error) string {
	fc.Check(err)
	return value
}

// Strings provides the []string part of the given fallible function result
// while registering the given error.
func (fc *FailureCollector) Strings(value []string, err error) []string {
	fc.Check(err)
	return value
}

// SyncStrategy provides the string part of the given fallible function result
// while registering the given error.
func (fc *FailureCollector) SyncStrategy(value config.SyncStrategy, err error) config.SyncStrategy {
	fc.Check(err)
	return value
}
