package gitdomain

// BranchesSnapshot is a snapshot of the Git branches at a particular point in time.
type BranchesSnapshot struct {
	// the branch that was checked out at the time the snapshot was taken
	Active LocalBranchName

	// Branches is a read-only copy of the branches that exist in this repo at the time the snapshot was taken.
	// Don't use these branches for business logic since businss logic might want to modify its in-memory cache of branches
	// as it adds or removes branches.
	Branches BranchInfos
}

func EmptyBranchesSnapshot() BranchesSnapshot {
	return BranchesSnapshot{
		Active:   EmptyLocalBranchName(),
		Branches: BranchInfos{},
	}
}

func (self BranchesSnapshot) IsEmpty() bool {
	return len(self.Branches) == 0 && self.Active.IsEmpty()
}
