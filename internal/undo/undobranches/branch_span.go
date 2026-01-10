package undobranches

import (
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/undo/undodomain"
	. "github.com/git-town/git-town/v22/pkg/prelude"
	"github.com/git-town/git-town/v22/pkg/set"
)

// BranchSpan represents changes of a branch over time.
type BranchSpan struct {
	Before Option[gitdomain.BranchInfo] // the status of the branch before Git Town ran
	After  Option[gitdomain.BranchInfo] // the status of the branch after Git Town ran
}

func (self BranchSpan) BranchNames() []gitdomain.BranchName {
	branchNames := set.New[gitdomain.BranchName]()
	if before, hasBefore := self.Before.Get(); hasBefore {
		if localName, hasLocalName := before.LocalName.Get(); hasLocalName {
			branchNames.Add(localName.BranchName())
		}
		if remoteName, hasRemoteName := before.RemoteName.Get(); hasRemoteName {
			branchNames.Add(remoteName.BranchName())
		}
	}
	if after, hasAfter := self.After.Get(); hasAfter {
		if localName, hasLocalName := after.LocalName.Get(); hasLocalName {
			branchNames.Add(localName.BranchName())
		}
		if remoteName, hasRemoteName := after.RemoteName.Get(); hasRemoteName {
			branchNames.Add(remoteName.BranchName())
		}
	}
	return branchNames.Values()
}

func (self BranchSpan) InconsistentChange() Option[undodomain.InconsistentChange] {
	_, isOmniChange := self.OmniChange().Get()
	localChanged, _, _, _ := self.LocalChanged()
	remoteChanged, _, _, _ := self.RemoteChanged()
	before, hasBefore := self.Before.Get()
	after, hasAfter := self.After.Get()
	isInconsistentChange := hasBefore && before.HasTrackingBranch() && hasAfter && after.HasTrackingBranch() && localChanged && remoteChanged && !isOmniChange
	if !isInconsistentChange {
		return None[undodomain.InconsistentChange]()
	}
	return Some(undodomain.InconsistentChange{
		After:  after,
		Before: before,
	})
}

// Indicates whether this BranchSpan describes the removal of an omni Branch
// and provides all relevant data for this situation.
func (self BranchSpan) IsOmniRemove() (isOmniRemove bool, beforeBranchName gitdomain.LocalBranchName, beforeSHA gitdomain.SHA) {
	before, hasBefore := self.Before.Get()
	if !hasBefore {
		return false, beforeBranchName, beforeSHA
	}
	beforeIsOmni, beforeName, beforeSHA := before.IsOmniBranch()
	_, hasAfter := self.After.Get()
	isOmniRemove = beforeIsOmni && !hasAfter
	return isOmniRemove, beforeName, beforeSHA
}

func (self BranchSpan) LocalAdded() (isLocalAdded bool, afterBranchName gitdomain.LocalBranchName, afterSHA gitdomain.SHA) {
	before, hasBefore := self.Before.Get()
	beforeHasLocalBranch, _, _ := before.GetLocal()
	after, hasAfter := self.After.Get()
	if !hasAfter {
		return false, afterBranchName, afterSHA
	}
	afterHasLocalBranch, afterLocalBranch, afterSHA := after.GetLocal()
	isLocalAdded = (!hasBefore || !beforeHasLocalBranch) && afterHasLocalBranch
	return isLocalAdded, afterLocalBranch, afterSHA
}

func (self BranchSpan) LocalChanged() (localChanged bool, branch gitdomain.LocalBranchName, beforeSHA, afterSHA gitdomain.SHA) {
	before, hasBefore := self.Before.Get()
	if !hasBefore {
		return false, branch, beforeSHA, afterSHA
	}
	hasLocalBranchBefore, beforeBranch, beforeSHA := before.GetLocal()
	after, hasAfter := self.After.Get()
	if !hasAfter {
		return false, branch, beforeSHA, afterSHA
	}
	hasLocalBranchAfter, _, afterSHA := after.GetLocal()
	localChanged = hasLocalBranchBefore && hasLocalBranchAfter && beforeSHA != afterSHA
	return localChanged, beforeBranch, beforeSHA, afterSHA
}

func (self BranchSpan) LocalRemoved() (localRemoved bool, branchName gitdomain.LocalBranchName, beforeSHA gitdomain.SHA) {
	before, hasBefore := self.Before.Get()
	hasBeforeBranch, branchName, beforeSHA := before.GetLocal()
	after, hasAfter := self.After.Get()
	hasAfterBranch, _, _ := after.GetLocal()
	localRemoved = hasBefore && hasBeforeBranch && (!hasAfter || !hasAfterBranch)
	return localRemoved, branchName, beforeSHA
}

// LocalRename indicates whether this BranchSpan describes the situation where only the local branch was renamed.
func (self BranchSpan) LocalRename() Option[LocalBranchRename] {
	before, hasBefore := self.Before.Get()
	if !hasBefore {
		return None[LocalBranchRename]()
	}
	after, hasAfter := self.After.Get()
	if !hasAfter {
		return None[LocalBranchRename]()
	}
	beforeName, hasBeforeName := before.LocalName.Get()
	if !hasBeforeName {
		return None[LocalBranchRename]()
	}
	afterName, hasAfterName := after.LocalName.Get()
	if !hasAfterName {
		return None[LocalBranchRename]()
	}
	beforeSHA, hasBeforeSHA := before.LocalSHA.Get()
	if !hasBeforeSHA {
		return None[LocalBranchRename]()
	}
	afterSHA, hasAfterSHA := after.LocalSHA.Get()
	if !hasAfterSHA {
		return None[LocalBranchRename]()
	}
	isLocalRename := beforeName != afterName && beforeSHA == afterSHA
	if !isLocalRename {
		return None[LocalBranchRename]()
	}
	return Some(LocalBranchRename{
		After:  afterName,
		Before: beforeName,
	})
}

// OmniChange indicates whether this BranchBeforeAfter changes a synced branch
// from one SHA both locally and remotely to another SHA both locally and remotely.
func (self BranchSpan) OmniChange() Option[LocalBranchChange] {
	before, hasBefore := self.Before.Get()
	if !hasBefore {
		return None[LocalBranchChange]()
	}
	beforeIsOmni, beforeName, beforeSHA := before.IsOmniBranch()
	after, hasAfter := self.After.Get()
	if !hasAfter {
		return None[LocalBranchChange]()
	}
	afterIsOmni, _, afterSHA := after.IsOmniBranch()
	isOmniChange := beforeIsOmni && afterIsOmni && beforeSHA != afterSHA
	if !isOmniChange {
		return None[LocalBranchChange]()
	}
	return Some(LocalBranchChange{
		beforeName: undodomain.Change[gitdomain.SHA]{
			Before: beforeSHA,
			After:  afterSHA,
		},
	})
}

func (self BranchSpan) RemoteAdded() (remoteAdded bool, addedRemoteBranchName gitdomain.RemoteBranchName, addedRemoteBranchSHA gitdomain.SHA) {
	before, hasBefore := self.Before.Get()
	beforeHasRemoteBranch, _, _ := before.GetRemote()
	after, hasAfter := self.After.Get()
	if !hasAfter {
		return false, addedRemoteBranchName, addedRemoteBranchSHA
	}
	afterHasRemoteBranch, afterRemoteBranchName, afterRemoteBranchSHA := after.GetRemote()
	remoteAdded = (!hasBefore || !beforeHasRemoteBranch) && afterHasRemoteBranch
	return remoteAdded, afterRemoteBranchName, afterRemoteBranchSHA
}

func (self BranchSpan) RemoteChanged() (remoteChanged bool, branchName gitdomain.RemoteBranchName, beforeSHA, afterSHA gitdomain.SHA) {
	before, hasBefore := self.Before.Get()
	if !hasBefore {
		return false, branchName, beforeSHA, afterSHA
	}
	beforeHasRemoteBranch, beforeRemoteBranchName, beforeRemoteBranchSHA := before.GetRemote()
	after, hasAfter := self.After.Get()
	if !hasAfter {
		return false, branchName, beforeSHA, afterSHA
	}
	afterHasRemoteBranch, _, afterRemoteBranchSHA := after.GetRemote()
	remoteChanged = beforeHasRemoteBranch && afterHasRemoteBranch && beforeRemoteBranchSHA != afterRemoteBranchSHA
	return remoteChanged, beforeRemoteBranchName, beforeRemoteBranchSHA, afterRemoteBranchSHA
}

func (self BranchSpan) RemoteRemoved() (remoteRemoved bool, remoteBranchName gitdomain.RemoteBranchName, beforeRemoteBranchSHA gitdomain.SHA) {
	before, hasBefore := self.Before.Get()
	if !hasBefore {
		return false, remoteBranchName, beforeRemoteBranchSHA
	}
	beforeHasRemoteBranch, remoteBranchName, beforeSHA := before.GetRemote()
	after, hasAfter := self.After.Get()
	afterHasRemoteBranch, _, _ := after.GetRemote()
	remoteRemoved = beforeHasRemoteBranch && (!hasAfter || !afterHasRemoteBranch)
	return remoteRemoved, remoteBranchName, beforeSHA
}

// func (self BranchSpan) String() string {
// 	result := "BranchSpan:\n"
// 	result += "Before:" + self.Before.String() + "\n"
// 	result += "After:" + self.After.String() + "\n"
// 	return result
// }
