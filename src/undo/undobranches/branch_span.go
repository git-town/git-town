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

// TODO: replace all uses with IsOmniRemove2
func (self BranchSpan) IsOmniRemove() bool {
	return self.Before.IsOmniBranch() && self.After.IsEmpty()
}

// Indicates whether this BranchSpan describes the removal of an omni Branch
// and provides all relevant data for this situation.
func (self BranchSpan) IsOmniRemove2() (isOmniRemove bool, beforeBranchName gitdomain.LocalBranchName, beforeSHA gitdomain.SHA) {
	beforeIsOmni, beforeName, beforeSHA := self.Before.IsOmni()
	isOmniRemove = beforeIsOmni && self.After.IsEmpty()
	return isOmniRemove, beforeName, beforeSHA
}

func (self BranchSpan) LocalAdded() bool {
	return !self.Before.HasLocalBranch() && self.After.HasLocalBranch()
}

func (self BranchSpan) LocalAdded2() (isLocalAdded bool, afterBranchName gitdomain.LocalBranchName, afterSHA gitdomain.SHA) {
	beforeHasLocalBranch, _, _ := self.Before.HasLocalBranch2()
	afterHasLocalBranch, afterLocalBranch, afterSHA := self.After.HasLocalBranch2()
	isLocalAdded = !beforeHasLocalBranch && afterHasLocalBranch
	return isLocalAdded, afterLocalBranch, afterSHA
}

func (self BranchSpan) LocalChanged() bool {
	return self.Before.LocalSHA != self.After.LocalSHA
}

func (self BranchSpan) LocalRemoved() bool {
	return self.Before.HasLocalBranch() && !self.After.HasLocalBranch()
}
func (self BranchSpan) LocalRemoved2() (localRemoved bool, beforeBranch gitdomain.LocalBranchName, beforeSHA gitdomain.SHA) {
	hasBeforeBranch, beforeBranch, beforeSHA := self.Before.HasLocalBranch2()
	hasAfterBranch, _, _ := self.After.HasLocalBranch2()
	localRemoved = hasBeforeBranch && !hasAfterBranch
	return localRemoved, beforeBranch, beforeSHA
}

// NoChanges indicates whether this BranchBeforeAfter contains changes or not.
func (self BranchSpan) NoChanges() bool {
	return !self.LocalChanged() && !self.RemoteChanged()
}

func (self BranchSpan) RemoteAdded() bool {
	return !self.Before.HasRemoteBranch() && self.After.HasRemoteBranch()
}

func (self BranchSpan) RemoteChanged() bool {
	return self.Before.RemoteSHA != self.After.RemoteSHA
}

func (self BranchSpan) RemoteRemoved() bool {
	return self.Before.HasRemoteBranch() && !self.After.HasRemoteBranch()
}
