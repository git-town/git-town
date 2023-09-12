package runstate

import "github.com/git-town/git-town/v9/src/domain"

// BranchesSnapshot is a snapshot of the Git branches at a particular point in time.
type BranchesSnapshot struct {

	// Branches is a read-only copy of the branches that exist in this repo at the time the snapshot was taken.
	// Don't use these branches for business logic since businss logic might want to modify its in-memory cache of branches
	// as it adds or removes branches.
	Branches domain.BranchInfos
}

func EmptyBranchesSnapshot() BranchesSnapshot {
	return BranchesSnapshot{}
}

func (bs BranchesSnapshot) Diff(other BranchesSnapshot) BranchesDiff {
	return BranchesDiff{}
}

type BranchesDiff struct {
	LocalAdded   domain.LocalBranchNames
	LocalRemoved map[domain.LocalBranchName]domain.SHA
	LocalChanged map[domain.LocalBranchName]Change[domain.SHA]
}

func (bd BranchesDiff) Steps() StepList {
	return StepList{}
}
