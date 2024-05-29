package undobranches

import (
	"github.com/git-town/git-town/v14/src/git/gitdomain"
	. "github.com/git-town/git-town/v14/src/gohacks/prelude"
)

// BranchSpan represents changes of a branch over time.
type BranchSpan struct {
	Before Option[gitdomain.BranchInfo] // the status of the branch before Git Town ran
	After  Option[gitdomain.BranchInfo] // the status of the branch after Git Town ran
}

func (self BranchSpan) IsInconsistentChange() (isInconsistentChange bool, before, after gitdomain.BranchInfo) {
	isOmniChange, _, _, _ := self.IsOmniChange()
	localChanged, _, _, _ := self.LocalChanged()
	remoteChanged, _, _, _ := self.RemoteChanged()
	before, hasBefore := self.Before.Get()
	after, hasAfter := self.After.Get()
	isInconsistentChange = hasBefore && before.HasTrackingBranch() && hasAfter && after.HasTrackingBranch() && localChanged && remoteChanged && !isOmniChange
	return isInconsistentChange, before, after
}

// IsOmniChange indicates whether this BranchBeforeAfter changes a synced branch
// from one SHA both locally and remotely to another SHA both locally and remotely.
func (self BranchSpan) IsOmniChange() (isOmniChange bool, branchName gitdomain.LocalBranchName, beforeSHA, afterSHA gitdomain.SHA) {
	before, hasBefore := self.Before.Get()
	if !hasBefore {
		return false, branchName, beforeSHA, afterSHA
	}
	beforeIsOmni, beforeName, beforeSHA := before.IsOmniBranch()
	after, hasAfter := self.After.Get()
	if !hasAfter {
		return false, branchName, beforeSHA, afterSHA
	}
	afterIsOmni, _, afterSHA := after.IsOmniBranch()
	isOmniChange = beforeIsOmni && afterIsOmni && beforeSHA != afterSHA
	return isOmniChange, beforeName, beforeSHA, afterSHA
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
	beforeHasLocalBranch, _, _ := before.HasLocalBranch()
	after, hasAfter := self.After.Get()
	if !hasAfter {
		return false, afterBranchName, afterSHA
	}
	afterHasLocalBranch, afterLocalBranch, afterSHA := after.HasLocalBranch()
	isLocalAdded = (!hasBefore || !beforeHasLocalBranch) && afterHasLocalBranch
	return isLocalAdded, afterLocalBranch, afterSHA
}

func (self BranchSpan) LocalChanged() (localChanged bool, branch gitdomain.LocalBranchName, beforeSHA, afterSHA gitdomain.SHA) {
	before, hasBefore := self.Before.Get()
	if !hasBefore {
		return false, branch, beforeSHA, afterSHA
	}
	hasLocalBranchBefore, beforeBranch, beforeSHA := before.HasLocalBranch()
	after, hasAfter := self.After.Get()
	if !hasAfter {
		return false, branch, beforeSHA, afterSHA
	}
	hasLocalBranchAfter, _, afterSHA := after.HasLocalBranch()
	localChanged = hasLocalBranchBefore && hasLocalBranchAfter && beforeSHA != afterSHA
	return localChanged, beforeBranch, beforeSHA, afterSHA
}

func (self BranchSpan) LocalRemoved() (localRemoved bool, branchName gitdomain.LocalBranchName, beforeSHA gitdomain.SHA) {
	before, hasBefore := self.Before.Get()
	hasBeforeBranch, branchName, beforeSHA := before.HasLocalBranch()
	after, hasAfter := self.After.Get()
	hasAfterBranch, _, _ := after.HasLocalBranch()
	localRemoved = hasBefore && hasBeforeBranch && (!hasAfter || !hasAfterBranch)
	return localRemoved, branchName, beforeSHA
}

// NoChanges indicates whether this BranchSpan contains changes or not.
func (self BranchSpan) NoChanges() bool {
	localAdded, _, _ := self.LocalAdded()
	localRemoved, _, _ := self.LocalRemoved()
	remoteAdded, _, _ := self.RemoteAdded()
	remoteRemoved, _, _ := self.RemoteRemoved()
	localChanged, _, _, _ := self.LocalChanged()
	remoteChanged, _, _, _ := self.RemoteChanged()
	return !localAdded && !localRemoved && !localChanged && !remoteAdded && !remoteRemoved && !remoteChanged
}

func (self BranchSpan) RemoteAdded() (remoteAdded bool, addedRemoteBranchName gitdomain.RemoteBranchName, addedRemoteBranchSHA gitdomain.SHA) {
	before, hasBefore := self.Before.Get()
	beforeHasRemoteBranch, _, _ := before.HasRemoteBranch()
	after, hasAfter := self.After.Get()
	if !hasAfter {
		return false, addedRemoteBranchName, addedRemoteBranchSHA
	}
	afterHasRemoteBranch, afterRemoteBranchName, afterRemoteBranchSHA := after.HasRemoteBranch()
	remoteAdded = (!hasBefore || !beforeHasRemoteBranch) && afterHasRemoteBranch
	return remoteAdded, afterRemoteBranchName, afterRemoteBranchSHA
}

func (self BranchSpan) RemoteChanged() (remoteChanged bool, branchName gitdomain.RemoteBranchName, beforeSHA, afterSHA gitdomain.SHA) {
	before, hasBefore := self.Before.Get()
	if !hasBefore {
		return false, branchName, beforeSHA, afterSHA
	}
	beforeHasRemoteBranch, beforeRemoteBranchName, beforeRemoteBranchSHA := before.HasRemoteBranch()
	after, hasAfter := self.After.Get()
	if !hasAfter {
		return false, branchName, beforeSHA, afterSHA
	}
	afterHasRemoteBranch, _, afterRemoteBranchSHA := after.HasRemoteBranch()
	remoteChanged = beforeHasRemoteBranch && afterHasRemoteBranch && beforeRemoteBranchSHA != afterRemoteBranchSHA
	return remoteChanged, beforeRemoteBranchName, beforeRemoteBranchSHA, afterRemoteBranchSHA
}

func (self BranchSpan) RemoteRemoved() (remoteRemoved bool, remoteBranchName gitdomain.RemoteBranchName, beforeRemoteBranchSHA gitdomain.SHA) {
	before, hasBefore := self.Before.Get()
	if !hasBefore {
		return false, remoteBranchName, beforeRemoteBranchSHA
	}
	beforeHasRemoteBranch, remoteBranchName, beforeSHA := before.HasRemoteBranch()
	after, hasAfter := self.After.Get()
	afterHasRemoteBranch, _, _ := after.HasRemoteBranch()
	remoteRemoved = beforeHasRemoteBranch && (!hasAfter || !afterHasRemoteBranch)
	return remoteRemoved, remoteBranchName, beforeSHA
}
