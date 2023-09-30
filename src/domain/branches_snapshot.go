package domain

// BranchesSnapshot is a snapshot of the Git branches at a particular point in time.
type BranchesSnapshot struct {
	// Branches is a read-only copy of the branches that exist in this repo at the time the snapshot was taken.
	// Don't use these branches for business logic since businss logic might want to modify its in-memory cache of branches
	// as it adds or removes branches.
	Branches BranchInfos

	// the branch that was checked out at the time the snapshot was taken
	Active LocalBranchName
}

func EmptyBranchesSnapshot() BranchesSnapshot {
	return BranchesSnapshot{
		Branches: BranchInfos{},
		Active:   LocalBranchName{},
	}
}

func (bs BranchesSnapshot) IsEmpty() bool {
	return len(bs.Branches) == 0 && bs.Active.IsEmpty()
}
