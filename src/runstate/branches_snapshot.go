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

// IsOmniChange indicates whether this BranchBeforeAfter changes a synced branch
// from one SHA both locally and remotely to another SHA both locally and remotely.
func (bba BranchBeforeAfter) IsOmniChange() bool {
	return bba.Before.IsOmniBranch() && bba.After.IsOmniBranch() && bba.LocalChanged()
}

func (bba BranchBeforeAfter) IsInconsintentChange() bool {
	return !bba.Before.LocalSHA.IsEmpty() &&
		!bba.Before.RemoteSHA.IsEmpty() &&
		!bba.After.LocalSHA.IsEmpty() &&
		!bba.After.RemoteSHA.IsEmpty() &&
		bba.LocalChanged() &&
		bba.RemoteChanged()
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
			result.OmniChanged[ba.Before.LocalName] = domain.Change[domain.SHA]{
				Before: ba.Before.LocalSHA,
				After:  ba.After.LocalSHA,
			}
			continue
		}
		if ba.IsInconsintentChange() {
			result.InconsistentlyChanged = append(result.InconsistentlyChanged, domain.InconsistentChange{
				Before: ba.Before,
				After:  ba.After,
			})
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
	RemoteRemoved domain.RemoteBranchesSHAs
	RemoteChanged domain.RemoteBranchChange
	// OmniChanges are changes where the local SHA and the remote SHA are identical before the change as well as after the change, and the SHA before and the SHA after are different.
	// Git Town recognizes OmniChanges because only they allow undoing changes made to remote perennial branches.
	// The reason is that perennial branches have protected remote branches, i.e. don't allow force-pushes to their remote branch. One can only do normal pushes.
	// So, to revert a change on a remote perennial branch one needs to perform a revert commit on the local perennial branch,
	// then normal-push (not force-push) that new commit up to the remote branch.
	// This is only possible if the local and remote branches have an identical SHA before as well as after.
	OmniChanged domain.LocalBranchChange // a branch had the same SHA locally and remotely, now it has a new SHA locally and remotely, the local and remote SHA are still equal
	// Inconsistent changes are changes on both local and tracking branch, but where the local and tracking branch don't have the same SHA before or after.
	// These changes cannot be undone for perennial branches because there is no way to reset the remote branch to the SHA it had before.
	InconsistentlyChanged domain.InconsistentChanges
}

// EmptyChanges provides a properly initialized empty Changes instance.
func EmptyChanges() Changes {
	return Changes{
		LocalAdded:            domain.LocalBranchNames{},
		LocalRemoved:          map[domain.LocalBranchName]domain.SHA{},
		LocalChanged:          domain.LocalBranchChange{},
		RemoteAdded:           []domain.RemoteBranchName{},
		RemoteRemoved:         map[domain.RemoteBranchName]domain.SHA{},
		RemoteChanged:         map[domain.RemoteBranchName]domain.Change[domain.SHA]{},
		OmniChanged:           domain.LocalBranchChange{},
		InconsistentlyChanged: domain.InconsistentChanges{},
	}
}

func (bd Changes) Steps(lineage config.Lineage, branchTypes domain.BranchTypes) StepList {
	result := StepList{}
	omniChangedPerennials, omniChangedFeatures := bd.OmniChanged.Categorize(branchTypes)

	// revert omni-changed perennial branches
	for branch, change := range omniChangedPerennials {
		result.Append(&steps.CheckoutStep{Branch: branch})
		result.Append(&steps.RevertCommitStep{SHA: change.Before})
		result.Append(&steps.PushCurrentBranchStep{CurrentBranch: branch, NoPushHook: false, Undoable: false})
	}

	// reset omni-changed feature branches
	for branch, change := range omniChangedFeatures {
		result.Append(&steps.CheckoutStep{Branch: branch})
		result.Append(&steps.ResetCurrentBranchToSHAStep{MustHaveSHA: change.After, SetToSHA: change.Before, Hard: true})
		result.Append(&steps.ForcePushBranchStep{Branch: branch, NoPushHook: false})
	}

	// ignore inconsistently changed perennial branches
	// because we can't change the remote and we therefore don't want to reset the local part either
	_, inconsistentChangedFeatures := bd.InconsistentlyChanged.Categorize(branchTypes)

	// reset inconsintently changed feature branches
	for _, inconsistentChange := range inconsistentChangedFeatures {
		result.Append(&steps.CheckoutStep{Branch: inconsistentChange.Before.LocalName})
		result.Append(&steps.ResetCurrentBranchToSHAStep{
			MustHaveSHA: inconsistentChange.After.LocalSHA,
			SetToSHA:    inconsistentChange.Before.LocalSHA,
			Hard:        true,
		})
		result.Append(&steps.ResetRemoteBranchToSHAStep{
			Branch:      inconsistentChange.Before.RemoteName,
			MustHaveSHA: inconsistentChange.After.RemoteSHA,
			SetToSHA:    inconsistentChange.Before.RemoteSHA,
		})
	}

	// remove remotely added branches
	for _, addedRemoteBranch := range bd.RemoteAdded {
		result.Append(&steps.DeleteTrackingBranchStep{
			Branch: addedRemoteBranch.LocalBranchName(),
		})
	}

	// re-create remotely removed feature branches
	_, removedFeatureTrackingBranches := bd.RemoteRemoved.Categorize(branchTypes)
	for branch, sha := range removedFeatureTrackingBranches {
		result.Append(&steps.CreateRemoteBranchStep{
			Branch:     branch.LocalBranchName(),
			SHA:        sha,
			NoPushHook: false,
		})
	}

	// reset locally changed branches
	for localBranch, change := range bd.LocalChanged {
		result.Append(&steps.CheckoutStep{Branch: localBranch})
		result.Append(&steps.ResetCurrentBranchToSHAStep{MustHaveSHA: change.After, SetToSHA: change.Before, Hard: true})
	}

	// remove locally added branches
	for _, addedLocalBranch := range bd.LocalAdded {
		result.Append(&steps.DeleteLocalBranchStep{
			Branch: addedLocalBranch,
			Parent: lineage.Parent(addedLocalBranch).Location(),
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

	_, remoteFeatureChanges := bd.RemoteChanged.Categorize(branchTypes)
	// Ignore remotely changed perennial branches because we can't force-push to them
	// and we would need the local branch to revert commits on them, but we can't change the local branch.

	// reset remotely changed feature branches
	for remoteChangedFeatureBranch, change := range remoteFeatureChanges {
		result.Append(&steps.ResetRemoteBranchToSHAStep{
			Branch:      remoteChangedFeatureBranch,
			MustHaveSHA: change.After,
			SetToSHA:    change.Before,
		})
	}
	return result
}
