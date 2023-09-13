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
	for _, before := range bs.Branches {
		if before.LocalName.IsEmpty() {
			// remote-only branch
			after := after.Branches.FindByRemote(before.RemoteName)
			if after == nil {
				result.RemoteRemoved[before.RemoteName] = before.RemoteSHA
				continue
			} else {
				// remote branch updated
				if before.RemoteSHA != after.RemoteSHA {
					result.RemoteChanged[before.RemoteName] = Change[domain.SHA]{
						Before: before.RemoteSHA,
						After:  after.RemoteSHA,
					}
					continue
				}
			}
		} else {
			// local branch
			after := after.Branches.FindLocalBranch(before.LocalName)
			if after == nil {
				result.LocalRemoved[before.LocalName] = before.LocalSHA
				continue
			}
			if before.LocalSHA != after.LocalSHA {
				result.LocalChanged[before.LocalName] = Change[domain.SHA]{
					Before: before.LocalSHA,
					After:  after.LocalSHA,
				}
				continue
			}
			if !before.RemoteSHA.IsEmpty() && after.RemoteSHA.IsEmpty() {
				result.RemoteRemoved[before.RemoteName] = before.RemoteSHA
				continue
			}
		}
	}
	for _, afterBranch := range after.Branches {
		if !afterBranch.LocalName.IsEmpty() {
			before := bs.Branches.FindLocalBranch(afterBranch.LocalName)
			if before == nil {
				result.LocalAdded = append(result.LocalAdded, afterBranch.LocalName)
				continue
			}
		}
		before := bs.Branches.FindByRemote(afterBranch.RemoteName)
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
