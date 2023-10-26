package undobranch

import "github.com/git-town/git-town/v9/src/domain"

// Span represents changes of a branch over time.
type Span struct {
	Before domain.BranchInfo // the status of the branch before Git Town ran
	After  domain.BranchInfo // the status of the branch after Git Town ran
}

func (self Span) IsInconsistentChange() bool {
	return self.Before.HasAllBranches() && self.After.HasAllBranches() && self.LocalChanged() && self.RemoteChanged() && !self.IsOmniChange()
}

// IsOmniChange indicates whether this BranchBeforeAfter changes a synced branch
// from one SHA both locally and remotely to another SHA both locally and remotely.
func (self Span) IsOmniChange() bool {
	return self.Before.IsOmniBranch() && self.After.IsOmniBranch() && self.LocalChanged()
}

func (self Span) IsOmniRemove() bool {
	return self.Before.IsOmniBranch() && self.After.IsEmpty()
}

func (self Span) LocalAdded() bool {
	return !self.Before.HasLocalBranch() && self.After.HasLocalBranch()
}

func (self Span) LocalChanged() bool {
	return self.Before.LocalSHA != self.After.LocalSHA
}

func (self Span) LocalRemoved() bool {
	return self.Before.HasLocalBranch() && !self.After.HasLocalBranch()
}

// NoChanges indicates whether this BranchBeforeAfter contains changes or not.
func (self Span) NoChanges() bool {
	return !self.LocalChanged() && !self.RemoteChanged()
}

func (self Span) RemoteAdded() bool {
	return !self.Before.HasRemoteBranch() && self.After.HasRemoteBranch()
}

func (self Span) RemoteChanged() bool {
	return self.Before.RemoteSHA != self.After.RemoteSHA
}

func (self Span) RemoteRemoved() bool {
	return self.Before.HasRemoteBranch() && !self.After.HasRemoteBranch()
}
