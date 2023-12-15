package configdomain

import (
	"fmt"

	"github.com/git-town/git-town/v11/src/domain"
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

// Bool provides the bool part of the given fallible function result
// while registering the given error.
func (self *FailureCollector) Bool(value bool, err error) bool {
	self.Check(err)
	return value
}

func (self *FailureCollector) Branches(value domain.Branches, err error) domain.Branches {
	self.Check(err)
	return value
}

func (self *FailureCollector) BranchesSyncStatus(value domain.BranchInfos, err error) domain.BranchInfos {
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

// Hosting provides the config.Hosting part of the given fallible function result
// while registering the given error.
func (self *FailureCollector) Hosting(value Hosting, err error) Hosting {
	self.Check(err)
	return value
}

func (self *FailureCollector) NewBranchPush(value NewBranchPush, err error) NewBranchPush {
	self.Check(err)
	return value
}

func (self *FailureCollector) Offline(value Offline, err error) Offline {
	self.Check(err)
	return value
}

func (self *FailureCollector) PushHook(value PushHook, err error) PushHook {
	self.Check(err)
	return value
}

// Remotes provides the domain.Remotes part of the given fallible function result
// while registering the given error.
func (self *FailureCollector) Remotes(value domain.Remotes, err error) domain.Remotes {
	self.Check(err)
	return value
}

func (self *FailureCollector) RepoStatus(value domain.RepoStatus, err error) domain.RepoStatus {
	self.Check(err)
	return value
}

func (self *FailureCollector) ShipDeleteRemoteBranch(value ShipDeleteTrackingBranch, err error) ShipDeleteTrackingBranch {
	self.Check(err)
	return value
}

// String provides the string part of the given fallible function result
// while registering the given error.
func (self *FailureCollector) String(value string, err error) string {
	self.Check(err)
	return value
}

// Strings provides the []string part of the given fallible function result
// while registering the given error.
func (self *FailureCollector) Strings(value []string, err error) []string {
	self.Check(err)
	return value
}

func (self *FailureCollector) SyncBeforeShip(value SyncBeforeShip, err error) SyncBeforeShip {
	self.Check(err)
	return value
}

// SyncFeatureStrategy provides the string part of the given fallible function result
// while registering the given error.
func (self *FailureCollector) SyncFeatureStrategy(value SyncFeatureStrategy, err error) SyncFeatureStrategy {
	self.Check(err)
	return value
}

// SyncPerennialStrategy provides the string part of the given fallible function result
// while registering the given error.
func (self *FailureCollector) SyncPerennialStrategy(value SyncPerennialStrategy, err error) SyncPerennialStrategy {
	self.Check(err)
	return value
}

func (self *FailureCollector) SyncUpstream(value SyncUpstream, err error) SyncUpstream {
	self.Check(err)
	return value
}
