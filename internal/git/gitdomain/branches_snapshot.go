package gitdomain

import . "github.com/git-town/git-town/v21/pkg/prelude"

// BranchesSnapshot is a snapshot of the Git branches at a particular point in time.
type BranchesSnapshot struct {
	// the branch that was checked out at the time the snapshot was taken
	Active Option[LocalBranchName]

	// Branches is a read-only copy of the branches that exist in this repo at the time the snapshot was taken.
	// Don't use these branches for business logic since businss logic might want to modify its in-memory cache of branches
	// as it adds or removes branches.
	Branches BranchInfos

	// Headless indicates whether the repo was in a headless state at the time this snapshot was taken.
	//
	// In this codebase we use the term "headless" instead of the more accurate term "detached head"
	// to distinguish this concept from Git Town's "detached" config setting and the "detach" command.
	Headless bool
}

func EmptyBranchesSnapshot() BranchesSnapshot {
	return BranchesSnapshot{
		Active:   None[LocalBranchName](),
		Branches: BranchInfos{},
		Headless: false,
	}
}
