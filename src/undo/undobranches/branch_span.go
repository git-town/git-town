package undobranches

import "github.com/git-town/git-town/v14/src/git/gitdomain"

// BranchSpan represents changes of a branch over time.
type BranchSpan struct {
	Before gitdomain.BranchInfo // the status of the branch before Git Town ran
	After  gitdomain.BranchInfo // the status of the branch after Git Town ran
}

func (self BranchSpan) IsInconsistentChange() bool {
	return self.Before.HasTrackingBranch() && self.After.HasTrackingBranch() && self.LocalChanged() && self.RemoteChanged() && !self.IsOmniChange()
}

// IsOmniChange indicates whether this BranchBeforeAfter changes a synced branch
// from one SHA both locally and remotely to another SHA both locally and remotely.
func (self BranchSpan) IsOmniChange() bool {
	return self.Before.IsOmniBranch() && self.After.IsOmniBranch() && self.LocalChanged()
}

// IsOmniChange indicates whether this BranchBeforeAfter changes a synced branch
// from one SHA both locally and remotely to another SHA both locally and remotely.
func (self BranchSpan) IsOmniChange2() (isOmniChange bool, branchName gitdomain.LocalBranchName, beforeSHA, afterSHA gitdomain.SHA) {
	beforeIsOmni, beforeName, beforeSHA := self.Before.IsOmni()
	afterIsOmni, _, afterSHA := self.After.IsOmni()
	isOmniChange = beforeIsOmni && afterIsOmni && beforeSHA != afterSHA
	return isOmniChange, beforeName, beforeSHA, afterSHA
}

// Indicates whether this BranchSpan describes the removal of an omni Branch
// and provides all relevant data for this situation.
func (self BranchSpan) IsOmniRemove() (isOmniRemove bool, beforeBranchName gitdomain.LocalBranchName, beforeSHA gitdomain.SHA) {
	beforeIsOmni, beforeName, beforeSHA := self.Before.IsOmni()
	isOmniRemove = beforeIsOmni && self.After.IsEmpty()
	return isOmniRemove, beforeName, beforeSHA
}

func (self BranchSpan) LocalAdded() (isLocalAdded bool, afterBranchName gitdomain.LocalBranchName, afterSHA gitdomain.SHA) {
	beforeHasLocalBranch, _, _ := self.Before.HasLocalBranch2()
	afterHasLocalBranch, afterLocalBranch, afterSHA := self.After.HasLocalBranch2()
	isLocalAdded = !beforeHasLocalBranch && afterHasLocalBranch
	return isLocalAdded, afterLocalBranch, afterSHA
}

func (self BranchSpan) LocalChanged() bool {
	return self.Before.LocalSHA != self.After.LocalSHA
}

func (self BranchSpan) LocalChanged2() (localChanged bool, branch gitdomain.LocalBranchName, beforeSHA, afterSHA gitdomain.SHA) {
	hasLocalBranchBefore, beforeBranch, beforeSHA := self.Before.HasLocalBranch2()
	hasLocalBranchAfter, _, afterSHA := self.After.HasLocalBranch2()
	localChanged = hasLocalBranchBefore && hasLocalBranchAfter && beforeSHA != afterSHA
	return localChanged, beforeBranch, beforeSHA, afterSHA
}

func (self BranchSpan) LocalRemoved() (localRemoved bool, branchName gitdomain.LocalBranchName, beforeSHA gitdomain.SHA) {
	hasBeforeBranch, branchName, beforeSHA := self.Before.HasLocalBranch2()
	hasAfterBranch, _, _ := self.After.HasLocalBranch2()
	localRemoved = hasBeforeBranch && !hasAfterBranch
	return localRemoved, branchName, beforeSHA
}

// NoChanges indicates whether this BranchBeforeAfter contains changes or not.
func (self BranchSpan) NoChanges() bool {
	return !self.LocalChanged() && !self.RemoteChanged()
}

func (self BranchSpan) RemoteAdded() bool {
	return !self.Before.HasRemoteBranch() && self.After.HasRemoteBranch()
}

func (self BranchSpan) RemoteAdded2() (remoteAdded bool, addedRemoteBranchName gitdomain.RemoteBranchName, addedRemoteBranchSHA gitdomain.SHA) {
	beforeHasRemoteBranch, _, _ := self.Before.HasRemoteBranch2()
	afterHasRemoteBranch, afterRemoteBranchName, afterRemoteBranchSHA := self.After.HasRemoteBranch2()
	remoteAdded = !beforeHasRemoteBranch && afterHasRemoteBranch
	return remoteAdded, afterRemoteBranchName, afterRemoteBranchSHA
}

func (self BranchSpan) RemoteChanged() bool {
	return self.Before.RemoteSHA != self.After.RemoteSHA
}

func (self BranchSpan) RemoteChanged2() (remoteChanged bool, branchName gitdomain.RemoteBranchName, beforeSHA, afterSHA gitdomain.SHA) {
	beforeHasRemoteBranch, beforeRemoteBranchName, beforeRemoteBranchSHA := self.Before.HasRemoteBranch2()
	afterHasRemoteBranch, _, afterRemoteBranchSHA := self.After.HasRemoteBranch2()
	remoteChanged = beforeHasRemoteBranch && afterHasRemoteBranch && beforeRemoteBranchSHA != afterRemoteBranchSHA
	return remoteChanged, beforeRemoteBranchName, beforeSHA, afterSHA
}

func (self BranchSpan) RemoteRemoved() bool {
	return self.Before.HasRemoteBranch() && !self.After.HasRemoteBranch()
}

func (self BranchSpan) RemoteRemoved2() (remoteRemoved bool, beforeRemoteBranchName gitdomain.RemoteBranchName, beforeRemoteBranchSHA gitdomain.SHA) {
	beforeHasRemoteBranch, beforeRemoteBranchName, beforeSHA := self.Before.HasRemoteBranch2()
	afterHasRemoteBranch, _, _ := self.After.HasRemoteBranch2()
	remoteRemoved = beforeHasRemoteBranch && !afterHasRemoteBranch
	return remoteRemoved, beforeRemoteBranchName, beforeSHA
}
