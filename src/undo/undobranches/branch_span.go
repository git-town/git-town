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

func (self BranchSpan) IsOmniRemove() bool {
	return self.Before.IsOmniBranch() && self.After.IsEmpty()
}

func (self BranchSpan) LocalAdded() bool {
	return !self.Before.HasLocalBranch() && self.After.HasLocalBranch()
}

func (self BranchSpan) LocalChanged() bool {
	return self.Before.LocalSHA != self.After.LocalSHA
}

func (self BranchSpan) LocalRemoved() bool {
	return self.Before.HasLocalBranch() && !self.After.HasLocalBranch()
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
