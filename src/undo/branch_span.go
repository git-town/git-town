package undo

import "github.com/git-town/git-town/v9/src/domain"

// BranchSpan represents changes of a branch over time.
type BranchSpan struct {
	Before domain.BranchInfo // the status of the branch before Git Town ran
	After  domain.BranchInfo // the status of the branch after Git Town ran
}

// IsOmniChange indicates whether this BranchBeforeAfter changes a synced branch
// from one SHA both locally and remotely to another SHA both locally and remotely.
func (bs BranchSpan) IsOmniChange() bool {
	return bs.Before.IsOmniBranch() && bs.After.IsOmniBranch() && bs.LocalChanged()
}

func (bs BranchSpan) IsOmniRemove() bool {
	return bs.Before.IsOmniBranch() && bs.After.IsEmpty()
}

func (bs BranchSpan) IsInconsistentChange() bool {
	return bs.Before.HasAllBranches() && bs.After.HasAllBranches() && bs.LocalChanged() && bs.RemoteChanged() && !bs.IsOmniChange()
}

func (bs BranchSpan) LocalAdded() bool {
	return !bs.Before.HasLocalBranch() && bs.After.HasLocalBranch()
}

func (bs BranchSpan) LocalChanged() bool {
	return bs.Before.LocalSHA != bs.After.LocalSHA
}

func (bs BranchSpan) LocalRemoved() bool {
	return bs.Before.HasLocalBranch() && !bs.After.HasLocalBranch()
}

// NoChanges indicates whether this BranchBeforeAfter contains changes or not.
func (bs BranchSpan) NoChanges() bool {
	return !bs.LocalChanged() && !bs.RemoteChanged()
}

func (bs BranchSpan) RemoteAdded() bool {
	return !bs.Before.HasRemoteBranch() && bs.After.HasRemoteBranch()
}

func (bs BranchSpan) RemoteChanged() bool {
	return bs.Before.RemoteSHA != bs.After.RemoteSHA
}

func (bs BranchSpan) RemoteRemoved() bool {
	return bs.Before.HasRemoteBranch() && !bs.After.HasRemoteBranch()
}
