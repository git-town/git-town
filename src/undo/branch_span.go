package undo

import "github.com/git-town/git-town/v9/src/domain"

// BranchSpan represents the temporal change of a branch.
type BranchSpan struct {
	Before domain.BranchInfo // the status of the branch before Git Town ran
	After  domain.BranchInfo // the status of the branch after Git Town ran
}

// IsOmniChange indicates whether this BranchBeforeAfter changes a synced branch
// from one SHA both locally and remotely to another SHA both locally and remotely.
func (bs BranchSpan) IsOmniChange() bool {
	return bs.Before.IsOmniBranch() && bs.After.IsOmniBranch() && bs.LocalChanged()
}

func (bba BranchSpan) IsInconsintentChange() bool {
	return !bba.Before.LocalSHA.IsEmpty() &&
		!bba.Before.RemoteSHA.IsEmpty() &&
		!bba.After.LocalSHA.IsEmpty() &&
		!bba.After.RemoteSHA.IsEmpty() &&
		bba.LocalChanged() &&
		bba.RemoteChanged()
}

func (b BranchSpan) LocalAdded() bool {
	return !b.Before.HasLocalBranch() && b.After.HasLocalBranch()
}

func (b BranchSpan) LocalChanged() bool {
	return b.Before.LocalSHA != b.After.LocalSHA
}

func (b BranchSpan) LocalRemoved() bool {
	return b.Before.HasLocalBranch() && !b.After.HasLocalBranch()
}

// NoChanges indicates whether this BranchBeforeAfter contains changes or not.
func (b BranchSpan) NoChanges() bool {
	return !b.LocalChanged() && !b.RemoteChanged()
}

func (b BranchSpan) RemoteAdded() bool {
	return !b.Before.HasRemoteBranch() && b.After.HasRemoteBranch()
}

func (b BranchSpan) RemoteChanged() bool {
	return b.Before.RemoteSHA != b.After.RemoteSHA
}

func (b BranchSpan) RemoteRemoved() bool {
	return b.Before.HasRemoteBranch() && !b.After.HasRemoteBranch()
}
