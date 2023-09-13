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

func (bs BranchesSnapshot) Diff(after BranchesSnapshot) BranchesDiff {
	result := BranchesDiff{
		LocalAdded:    domain.LocalBranchNames{},
		LocalRemoved:  map[domain.LocalBranchName]domain.SHA{},
		LocalChanged:  map[domain.LocalBranchName]Change[domain.SHA]{},
		RemoteAdded:   []domain.RemoteBranchName{},
		RemoteRemoved: map[domain.RemoteBranchName]domain.SHA{},
		RemoteChanged: map[domain.RemoteBranchName]Change[domain.SHA]{},
	}
	for _, beforeBranch := range bs.Branches {
		afterBI := after.Branches.FindLocalBranch(beforeBranch.LocalName)
		if afterBI == nil {
			result.LocalRemoved[beforeBranch.LocalName] = beforeBranch.LocalSHA
			continue
		}
		if beforeBranch.LocalSHA != afterBI.LocalSHA {
			result.LocalChanged[beforeBranch.LocalName] = Change[domain.SHA]{
				Before: beforeBranch.LocalSHA,
				After:  afterBI.LocalSHA,
			}
		}
		if beforeBranch.RemoteSHA != afterBI.RemoteSHA {

		}
	}
	for _, afterBranch := range after.Branches {
		before := bs.Branches.FindLocalBranch(afterBranch.LocalName)
		if before == nil {
			result.LocalAdded = append(result.LocalAdded, afterBranch.LocalName)
			continue
		}
		before = bs.Branches.FindLocalBranchWithTracking(afterBranch.RemoteName)
		if before == nil {
			result.RemoteAdded = append(result.RemoteAdded, afterBranch.RemoteName)
			continue
		}
	}
	return result
}

type BranchesDiff struct {
	LocalAdded    domain.LocalBranchNames
	LocalRemoved  map[domain.LocalBranchName]domain.SHA
	LocalChanged  map[domain.LocalBranchName]Change[domain.SHA]
	RemoteAdded   []domain.RemoteBranchName
	RemoteRemoved map[domain.RemoteBranchName]domain.SHA
	RemoteChanged map[domain.RemoteBranchName]Change[domain.SHA]
}

func (bd BranchesDiff) Steps() StepList {
	return StepList{}
}
