package runstate

import (
	"github.com/git-town/git-town/v9/src/config"
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/steps"
)

// BranchesSnapshot is a snapshot of the Git branches at a particular point in time.
type BranchesSnapshot struct {
	// Branches is a read-only copy of the branches that exist in this repo at the time the snapshot was taken.
	// Don't use these branches for business logic since businss logic might want to modify its in-memory cache of branches
	// as it adds or removes branches.
	Branches domain.BranchInfos
}

func EmptyBranchesSnapshot() BranchesSnapshot {
	return BranchesSnapshot{
		Branches: domain.BranchInfos{},
	}
}

// BranchBeforeAfter represents the temporal change of a branch.
type BranchBeforeAfter struct {
	Before domain.BranchInfo // the status of the branch before Git Town ran
	After  domain.BranchInfo // the status of the branch after Git Town ran
}

// IsOmniAdd indicates whether this BranchBeforeAfter adds an omnibranch.
func (bba BranchBeforeAfter) IsOmniAdd() bool {
	return bba.Before.IsEmpty() && !bba.After.IsEmpty() && bba.After.IsOmniBranch()
}

// IsOmniChange indicates whether this BranchBeforeAfter changes a synced branch
// from one SHA both locally and remotely to another SHA both locally and remotely.
func (bba BranchBeforeAfter) IsOmniChange() bool {
	return bba.Before.IsOmniBranch() && bba.After.IsOmniBranch() && bba.LocalChanged()
}

// IsOmniRemove indicates whether this BranchBeforeAfter removes an omnibranch.
func (bba BranchBeforeAfter) IsOmniRemove() bool {
	return !bba.Before.IsEmpty() && bba.Before.IsOmniBranch() && bba.After.IsEmpty()
}

func (bba BranchBeforeAfter) LocalAdded() bool {
	return !bba.Before.HasLocalBranch() && bba.After.HasLocalBranch()
}

func (bba BranchBeforeAfter) LocalChanged() bool {
	return bba.Before.LocalSHA != bba.After.LocalSHA
}

func (bba BranchBeforeAfter) LocalRemoved() bool {
	return bba.Before.HasLocalBranch() && !bba.After.HasLocalBranch()
}

// NoChanges indicates whether this BranchBeforeAfter contains changes or not.
func (bba BranchBeforeAfter) NoChanges() bool {
	return !bba.LocalChanged() && !bba.RemoteChanged()
}

func (bba BranchBeforeAfter) RemoteAdded() bool {
	return !bba.Before.HasRemoteBranch() && bba.After.HasRemoteBranch()
}

func (bba BranchBeforeAfter) RemoteChanged() bool {
	return bba.Before.RemoteSHA != bba.After.RemoteSHA
}

func (bba BranchBeforeAfter) RemoteRemoved() bool {
	return bba.Before.HasRemoteBranch() && !bba.After.HasRemoteBranch()
}

type BranchesBeforeAfter []BranchBeforeAfter

func (bs BranchesSnapshot) Changes(afterSnapshot BranchesSnapshot) BranchesBeforeAfter {
	result := BranchesBeforeAfter{}
	for _, before := range bs.Branches {
		after := afterSnapshot.Branches.FindMatchingRecord(before)
		result = append(result, BranchBeforeAfter{
			Before: before,
			After:  after,
		})
	}
	for _, after := range afterSnapshot.Branches {
		if bs.Branches.FindMatchingRecord(after).IsEmpty() {
			result = append(result, BranchBeforeAfter{
				Before: domain.EmptyBranchInfo(),
				After:  after,
			})
		}
	}
	return result
}

// Diff describes the changes made in this BranchesBeforeAfter structure.
func (bc BranchesBeforeAfter) Diff() Changes {
	result := EmptyChanges()
	for _, ba := range bc {
		if ba.NoChanges() {
			continue
		}
		if ba.IsOmniChange() {
			result.BothChanged[ba.Before.LocalName] = domain.Change[domain.SHA]{
				Before: ba.Before.LocalSHA,
				After:  ba.After.LocalSHA,
			}
			continue
		}
		if ba.IsOmniAdd() {
			result.BothAdded = append(result.BothAdded, ba.After.LocalName)
			continue
		}
		if ba.IsOmniRemove() {
			result.BothRemoved[ba.Before.LocalName] = ba.Before.LocalSHA
			continue
		}
		switch {
		case ba.LocalAdded():
			result.LocalAdded = append(result.LocalAdded, ba.After.LocalName)
		case ba.LocalRemoved():
			result.LocalRemoved[ba.Before.LocalName] = ba.Before.LocalSHA
		case ba.LocalChanged():
			result.LocalChanged[ba.Before.LocalName] = domain.Change[domain.SHA]{
				Before: ba.Before.LocalSHA,
				After:  ba.After.LocalSHA,
			}
		}
		switch {
		case ba.RemoteAdded():
			result.RemoteAdded = append(result.RemoteAdded, ba.After.RemoteName)
		case ba.RemoteRemoved():
			result.RemoteRemoved[ba.Before.RemoteName] = ba.Before.RemoteSHA
		case ba.RemoteChanged():
			result.RemoteChanged[ba.Before.RemoteName] = domain.Change[domain.SHA]{
				Before: ba.Before.RemoteSHA,
				After:  ba.After.RemoteSHA,
			}
		}
	}
	return result
}

type Changes struct {
	LocalAdded    domain.LocalBranchNames
	LocalRemoved  map[domain.LocalBranchName]domain.SHA
	LocalChanged  domain.LocalBranchChange
	RemoteAdded   []domain.RemoteBranchName
	RemoteRemoved map[domain.RemoteBranchName]domain.SHA
	RemoteChanged map[domain.RemoteBranchName]domain.Change[domain.SHA]
	BothAdded     domain.LocalBranchNames               // a branch was added locally and remotely with the same SHA
	BothRemoved   map[domain.LocalBranchName]domain.SHA // a branch that had the same SHA locally and remotely was removed from both locations
	BothChanged   domain.LocalBranchChange              // a branch that had the same SHA locally and remotely now has a new SHA locally and remotely, and the local and remote SHA are still equal
}

// EmptyChanges provides a properly initialized empty Changes instance.
func EmptyChanges() Changes {
	return Changes{
		LocalAdded:    domain.LocalBranchNames{},
		LocalRemoved:  map[domain.LocalBranchName]domain.SHA{},
		LocalChanged:  domain.LocalBranchChange{},
		RemoteAdded:   []domain.RemoteBranchName{},
		RemoteRemoved: map[domain.RemoteBranchName]domain.SHA{},
		RemoteChanged: map[domain.RemoteBranchName]domain.Change[domain.SHA]{},
		BothAdded:     domain.LocalBranchNames{},
		BothRemoved:   map[domain.LocalBranchName]domain.SHA{},
		BothChanged:   domain.LocalBranchChange{},
	}
}

func (bd Changes) Steps(lineage config.Lineage, branchTypes domain.BranchTypes) StepList {
	result := StepList{}
	bothChangedPerennials, bothChangedFeatures := bd.BothChanged.Categorize(branchTypes)
	// revert omni-changed perennial branches
	for branch, change := range bothChangedPerennials {
		result.Append(&steps.CheckoutStep{Branch: branch})
		result.Append(&steps.RevertCommitStep{SHA: change.Before})
		result.Append(&steps.ForcePushBranchStep{Branch: branch, NoPushHook: false})
	}
	// reset omni-changed feature branches
	for branch, change := range bothChangedFeatures {
		result.Append(&steps.CheckoutStep{Branch: branch})
		result.Append(&steps.ResetCurrentBranchToSHAStep{SHA: change.Before, Hard: true})
		result.Append(&steps.ForcePushBranchStep{Branch: branch, NoPushHook: false})
	}

	// delete omni-changed feature branches locally and push updates
	// remove added omni-branches
	// re-create removed omni-branches

	// revert locally changed perennial branches
	for localBranch, change := range bd.LocalChanged {
		result.Append(&steps.CheckoutStep{Branch: localBranch})
		result.Append(&steps.ResetCurrentBranchToSHAStep{SHA: change.Before, Hard: true})
	}

	// reset locally changed feature branches
	// remove locally added branches
	for _, addedLocalBranch := range bd.LocalAdded {
		parent := lineage.Parent(addedLocalBranch)
		result.Append(&steps.DeleteLocalBranchStep{
			Branch: addedLocalBranch,
			Parent: parent.Location(),
			Force:  true,
		})
	}
	// re-create locally removed branches
	for removedLocalBranch, startingPoint := range bd.LocalRemoved {
		result.Append(&steps.CreateBranchStep{
			Branch:        removedLocalBranch,
			StartingPoint: startingPoint.Location(),
		})
	}

	// revert remotely changed perennial branches
	// reset remotely changed feature branches
	// remove remotely added branches
	// re-create remotely removed branches
	return result
}
