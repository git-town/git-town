package domain

import (
	"fmt"

	"github.com/git-town/git-town/v9/src/messages"
)

// BranchInfos contains the BranchInfos for all branches in a repo.
// Tracking branches on the origin remote don't get their own entry,
// they are listed in the `TrackingBranch` property of the local branch they track.
type BranchInfos []BranchInfo

func (bis BranchInfos) Clone() BranchInfos {
	// appending to a slice with zero capacity (zero value) allocates only once
	return append(bis[:0:0], bis...)
}

// FindByLocalName provides the branch with the given name if one exists.
func (bis BranchInfos) FindByLocalName(branchName LocalBranchName) *BranchInfo {
	for bi, branch := range bis {
		if branch.LocalName == branchName {
			return &bis[bi]
		}
	}
	return nil
}

// FindByRemoteName provides the local branch that has the given remote branch as its tracking branch
// or nil if no such branch exists.
func (bis BranchInfos) FindByRemoteName(remoteBranch RemoteBranchName) *BranchInfo {
	for b, bi := range bis {
		if bi.RemoteName == remoteBranch {
			return &bis[b]
		}
	}
	return nil
}

func (bis BranchInfos) FindMatchingRecord(other BranchInfo) BranchInfo {
	for _, bi := range bis {
		if bi.LocalName == other.LocalName && !other.LocalName.IsEmpty() {
			return bi
		}
		if bi.RemoteName == other.RemoteName && !other.RemoteName.IsEmpty() {
			return bi
		}
	}
	return EmptyBranchInfo()
}

// IsKnown indicates whether the given local branch is already known to this BranchesSyncStatus instance.
func (bis BranchInfos) HasLocalBranch(localBranch LocalBranchName) bool {
	for _, bi := range bis {
		if bi.LocalName == localBranch {
			return true
		}
	}
	return false
}

// HasMatchingRemoteBranchFor indicates whether there is already a remote branch matching the given local branch.
func (bis BranchInfos) HasMatchingRemoteBranchFor(localBranch LocalBranchName) bool {
	// TODO: rename .RemoteBranch to .TrackingBranchName
	return bis.FindByRemoteName(localBranch.RemoteBranch()) != nil
}

// LocalBranches provides only the branches that exist on the local machine.
func (bis BranchInfos) LocalBranches() BranchInfos {
	result := BranchInfos{}
	for _, bi := range bis {
		if bi.IsLocal() {
			result = append(result, bi)
		}
	}
	return result
}

// LocalBranchesWithDeletedTrackingBranches provides only the branches that exist locally and have a deleted tracking branch.
func (bis BranchInfos) LocalBranchesWithDeletedTrackingBranches() BranchInfos {
	result := BranchInfos{}
	for _, bi := range bis {
		if bi.SyncStatus == SyncStatusDeletedAtRemote {
			result = append(result, bi)
		}
	}
	return result
}

// Names provides the names of all local branches in this BranchesSyncStatus instance.
func (bis BranchInfos) Names() LocalBranchNames {
	result := make(LocalBranchNames, 0, len(bis))
	for _, bi := range bis {
		if !bi.LocalName.IsEmpty() {
			result = append(result, bi.LocalName)
		}
	}
	return result
}

func (bis BranchInfos) Remove(branchName LocalBranchName) BranchInfos {
	result := BranchInfos{}
	for _, bi := range bis {
		if bi.LocalName != branchName {
			result = append(result, bi)
		}
	}
	return result
}

// Select provides the BranchSyncStatus elements with the given names.
func (bis BranchInfos) Select(names []LocalBranchName) (BranchInfos, error) {
	result := make(BranchInfos, len(names))
	for b, bi := range names {
		branch := bis.FindByLocalName(bi)
		if branch == nil {
			return result, fmt.Errorf(messages.BranchDoesntExist, bi)
		}
		result[b] = *branch
	}
	return result, nil
}

func (bis BranchInfos) UpdateLocalSHA(branch LocalBranchName, sha SHA) error {
	for b := range bis {
		if bis[b].LocalName == branch {
			bis[b].LocalSHA = sha
			return nil
		}
	}
	return fmt.Errorf("branch %q not found", branch)
}
