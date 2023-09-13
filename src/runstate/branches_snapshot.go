package runstate

import (
	"github.com/git-town/git-town/v9/src/domain"
)

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

// BranchBeforeAfter represents the temporal change of a branch.
type BranchBeforeAfter struct {
	Before domain.BranchInfo // the status of the branch before Git Town ran
	After  domain.BranchInfo // the status of the branch after Git Town ran
}

// NoChanges indicates whether this BranchBeforeAfter contains changes or not.
func (bba BranchBeforeAfter) NoChanges() bool {
	return bba.Before.LocalSHA == bba.After.LocalSHA && bba.Before.RemoteSHA == bba.After.RemoteSHA
}

// IsOmniChange indicates whether this BranchBeforeAfter changes a synced branch
// from one SHA both locally and remotely to another SHA both locally and remotely.
func (bba BranchBeforeAfter) IsOmniChange() bool {
	if bba.Before.LocalSHA != bba.Before.RemoteSHA {
		// was not an omnibranch before
		return false
	}
	if bba.After.LocalSHA != bba.After.RemoteSHA {
		// is not an omnibranch after
		return false
	}
	if bba.Before.LocalSHA == bba.After.LocalSHA {
		// no change
		return false
	}
	return true
}

type BranchesBeforeAfter []BranchBeforeAfter

// Changes provides matching BranchInfos from before and after.
func (bs BranchesSnapshot) Changes(afterSnapshot BranchesSnapshot) BranchesBeforeAfter {
	result := BranchesBeforeAfter{}
	// for _, before := range bs.Branches {
	// 	if !before.LocalName.IsEmpty() {
	// 		after := afterSnapshot.Branches.FindLocalBranch(before.LocalName)
	// 		if after != nil {
	// 			result = append(result, BranchBeforeAfter{
	// 				Before: &before,
	// 				After:  after,
	// 			})
	// 			continue
	// 		}
	// 	}
	// 	if !before.RemoteName.IsEmpty() {
	// 		after := afterSnapshot.Branches.FindByRemote(before.RemoteName)
	// 		result = append(result, BranchBeforeAfter{
	// 			Before: &before,
	// 			After:  after,
	// 		})
	// 	}
	// }
	// for _, after := range afterSnapshot.Branches {
	// 	if !after.LocalName.IsEmpty() {
	// 		before := bs.Branches.FindLocalBranch(after.LocalName)
	// 		if before == nil {
	// 			result = append(result, BranchBeforeAfter{
	// 				Before: nil,
	// 				After:  &after,
	// 			})
	// 		} else {
	// 			// here there exists a matching before and after --> it was already added when iterating before
	// 		}
	// 		continue
	// 	}
	// 	if !after.RemoteName.IsEmpty() {
	// 		before := bs.Branches.FindByRemote(after.RemoteName)
	// 		if before == nil {
	// 			result = append(result, BranchBeforeAfter{
	// 				Before: nil,
	// 				After:  &after,
	// 			})
	// 		} else {
	// 			// here there exists a matching before and after --> it was already added when iterating before
	// 		}
	// 	}
	// }
	return result
}

// Diff describes the changes made in this BranchesBeforeAfter structure.
func (bc BranchesBeforeAfter) Diff() Changes {
	result := Changes{
		LocalAdded:    domain.LocalBranchNames{},
		LocalRemoved:  map[domain.LocalBranchName]domain.SHA{},
		LocalChanged:  map[domain.LocalBranchName]Change[domain.SHA]{},
		RemoteAdded:   []domain.RemoteBranchName{},
		RemoteRemoved: map[domain.RemoteBranchName]domain.SHA{},
		RemoteChanged: map[domain.RemoteBranchName]Change[domain.SHA]{},
		BothAdded:     domain.LocalBranchNames{},
		BothRemoved:   map[domain.LocalBranchName]domain.SHA{},
		BothChanged:   map[domain.LocalBranchName]Change[domain.SHA]{},
	}
	// for _, ba := range bc {
	// 	if ba.NoChanges() {
	// 		continue
	// 	}
	// 	// check for omnibranch change
	// 	if beforeLocalSHA == beforeRemoteSHA && afterLocalSHA == afterRemoteSHA && beforeLocalSHA != afterLocalSHA {
	// 		result.BothChanged[ba.Before.LocalName] = Change[domain.SHA]{
	// 			Before: beforeLocalSHA,
	// 			After: afterLocalSHA,
	// 		}
	// 		continue
	// 	}
	// 	// check for omnibranch added
	// 	if beforeLocalSHA.IsEmpty() && beforeRemoteSHA.IsEmpty && !afterLocalSHA.IsEmpty() && afterLocalSHA == afterRemoteSHA {
	// 		result.BothAdded = append(result.BothAdded, ba.After.LocalName)
	// 		continue
	// 	}
	// 	// check for omnibranch removed
	// 	if !beforeLocalSHA.IsEmpty() && beforeLocalSHA == beforeRemoteSHA && afterLocalSHA.IsEmpty() && afterRemoteSHA.IsEmpty() {
	// 		result.BothRemoved[ba.Before.LocalName] = beforeLocalSHA
	// 		continue
	// 	}
	// 	if

	// 	if ba.Before != nil && ba.After != nil {
	// 		if !ba.Before.LocalName.IsEmpty() && !ba.After.LocalName.IsEmpty() {
	// 			if ba.Before.LocalSHA == ba.After.LocalSHA
	// 		}
	// 	}
	// 	if ba.Before != nil && ba.After == nil {
	// 	}
	// 	if ba.Before == nil && ba.After != nil {
	// 	}
	// 	panic("before and after are nil, this should never happen")
	// }

	// for _, before := range bs.Branches {
	// 	if before.LocalName.IsEmpty() {
	// 		// remote-only branch
	// 		after := after.Branches.FindByRemote(before.RemoteName)
	// 		if after == nil {
	// 			result.RemoteRemoved[before.RemoteName] = before.RemoteSHA
	// 			continue
	// 		} else {
	// 			// remote branch updated
	// 			if before.RemoteSHA != after.RemoteSHA {
	// 				result.RemoteChanged[before.RemoteName] = Change[domain.SHA]{
	// 					Before: before.RemoteSHA,
	// 					After:  after.RemoteSHA,
	// 				}
	// 				continue
	// 			}
	// 		}
	// 	} else {
	// 		// local or omni branch
	// 		after := after.Branches.FindLocalBranch(before.LocalName)
	// 		if after == nil {
	// 			result.LocalRemoved[before.LocalName] = before.LocalSHA
	// 			continue
	// 		}
	// 		if before.LocalSHA != after.LocalSHA {
	// 			result.LocalChanged[before.LocalName] = Change[domain.SHA]{
	// 				Before: before.LocalSHA,
	// 				After:  after.LocalSHA,
	// 			}
	// 			continue
	// 		}
	// 		if !before.RemoteSHA.IsEmpty() && after.RemoteSHA.IsEmpty() {
	// 			result.RemoteRemoved[before.RemoteName] = before.RemoteSHA
	// 			continue
	// 		}
	// 	}
	// }
	// for _, afterBranch := range after.Branches {
	// 	if !afterBranch.LocalName.IsEmpty() {
	// 		before := bs.Branches.FindLocalBranch(afterBranch.LocalName)
	// 		if before == nil {
	// 			result.LocalAdded = append(result.LocalAdded, afterBranch.LocalName)
	// 			continue
	// 		}
	// 	}
	// 	before := bs.Branches.FindByRemote(afterBranch.RemoteName)
	// 	if before == nil {
	// 		result.RemoteAdded = append(result.RemoteAdded, afterBranch.RemoteName)
	// 		continue
	// 	}
	// }
	return result
}

type Changes struct {
	LocalAdded    domain.LocalBranchNames
	LocalRemoved  map[domain.LocalBranchName]domain.SHA
	LocalChanged  map[domain.LocalBranchName]Change[domain.SHA]
	RemoteAdded   []domain.RemoteBranchName
	RemoteRemoved map[domain.RemoteBranchName]domain.SHA
	RemoteChanged map[domain.RemoteBranchName]Change[domain.SHA]
	BothAdded     domain.LocalBranchNames                       // a branch was added locally and remotely with the same SHA
	BothRemoved   map[domain.LocalBranchName]domain.SHA         // a branch that had the same SHA locally and remotely was removed from both locations
	BothChanged   map[domain.LocalBranchName]Change[domain.SHA] // a branch that had the same SHA locally and remotely now has a new SHA locally and remotely, and the local and remote SHA are still equal
}

func (bd Changes) Steps() StepList {
	return StepList{}
}
