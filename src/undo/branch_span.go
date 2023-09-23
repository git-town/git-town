package undo

import "github.com/git-town/git-town/v9/src/domain"

// BranchSpan represents the states of a branch before and after a change.
type BranchSpan struct {
	Before domain.BranchInfo // the status of the branch before Git Town ran
	After  domain.BranchInfo // the status of the branch after Git Town ran
}

// IsOmniChange indicates whether this BranchBeforeAfter changes a synced branch
// from one SHA both locally and remotely to another SHA both locally and remotely.
func (b BranchSpan) IsOmniChange() bool {
	return b.Before.IsOmniBranch() && b.After.IsOmniBranch() && b.LocalChanged()
}

func (b BranchSpan) IsOmniRemove() bool {
	return b.Before.IsOmniBranch() && b.After.IsEmpty()
}

func (b BranchSpan) IsInconsistentChange() bool {
	return !b.Before.LocalSHA.IsEmpty() &&
		!b.Before.RemoteSHA.IsEmpty() &&
		!b.After.LocalSHA.IsEmpty() &&
		!b.After.RemoteSHA.IsEmpty() &&
		b.LocalChanged() &&
		b.RemoteChanged()
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
