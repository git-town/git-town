package domain

import (
	"fmt"

	"github.com/git-town/git-town/v9/src/messages"
)

// BranchInfos contains the BranchInfos for all branches in a repo.
// Tracking branches on the origin remote don't get their own entry,
// they are listed in the `TrackingBranch` property of the local branch they track.
type BranchInfos []BranchInfo

func (bs BranchInfos) Clone() BranchInfos {
	// appending to a slice with zero capacity (zero value) allocates only once
	return append(bs[:0:0], bs...)
}

// FindByLocalName provides the branch with the given name if one exists.
func (bs BranchInfos) FindByLocalName(branchName LocalBranchName) *BranchInfo {
	for bi, branch := range bs {
		if branch.LocalName == branchName {
			return &bs[bi]
		}
	}
	return nil
}

// FindByRemoteName provides the local branch that has the given remote branch as its tracking branch
// or nil if no such branch exists.
func (bs BranchInfos) FindByRemoteName(remoteBranch RemoteBranchName) *BranchInfo {
	for b, branch := range bs {
		if branch.RemoteName == remoteBranch {
			return &bs[b]
		}
	}
	return nil
}

func (bs BranchInfos) FindMatchingRecord(other BranchInfo) BranchInfo {
	for _, branchInfo := range bs {
		if branchInfo.LocalName == other.LocalName && !other.LocalName.IsEmpty() {
			return branchInfo
		}
		if branchInfo.RemoteName == other.RemoteName && !other.RemoteName.IsEmpty() {
			return branchInfo
		}
	}
	return EmptyBranchInfo()
}

// IsKnown indicates whether the given local branch is already known to this BranchesSyncStatus instance.
func (bs BranchInfos) HasLocalBranch(localBranch LocalBranchName) bool {
	for _, branch := range bs {
		if branch.LocalName == localBranch {
			return true
		}
	}
	return false
}

// HasMatchingTrackingBranchFor indicates whether there is already a remote branch tracking the given local branch.
func (bs BranchInfos) HasMatchingTrackingBranchFor(localBranch LocalBranchName) bool {
	return bs.FindByRemoteName(localBranch.TrackingBranch()) != nil
}

// LocalBranches provides only the branches that exist on the local machine.
func (bs BranchInfos) LocalBranches() BranchInfos {
	result := BranchInfos{}
	for _, branch := range bs {
		if branch.IsLocal() {
			result = append(result, branch)
		}
	}
	return result
}

// LocalBranchesWithDeletedTrackingBranches provides only the branches that exist locally and have a deleted tracking branch.
func (bs BranchInfos) LocalBranchesWithDeletedTrackingBranches() BranchInfos {
	result := BranchInfos{}
	for _, branch := range bs {
		if branch.SyncStatus == SyncStatusDeletedAtRemote {
			result = append(result, branch)
		}
	}
	return result
}

// Names provides the names of all local branches in this BranchesSyncStatus instance.
func (bs BranchInfos) Names() LocalBranchNames {
	result := make(LocalBranchNames, 0, len(bs))
	for _, branch := range bs {
		if !branch.LocalName.IsEmpty() {
			result = append(result, branch.LocalName)
		}
	}
	return result
}

func (bs BranchInfos) Remove(branchName LocalBranchName) BranchInfos {
	result := BranchInfos{}
	for _, branch := range bs {
		if branch.LocalName != branchName {
			result = append(result, branch)
		}
	}
	return result
}

// Select provides the BranchSyncStatus elements with the given names.
func (bs BranchInfos) Select(names []LocalBranchName) (BranchInfos, error) {
	result := make(BranchInfos, len(names))
	for n, name := range names {
		branch := bs.FindByLocalName(name)
		if branch == nil {
			return result, fmt.Errorf(messages.BranchDoesntExist, name)
		}
		result[n] = *branch
	}
	return result, nil
}

// TODO: rename bs to bis.
func (bs BranchInfos) UpdateLocalSHA(branch LocalBranchName, sha SHA) error {
	for b := range bs {
		if bs[b].LocalName == branch {
			bs[b].LocalSHA = sha
			return nil
		}
	}
	return fmt.Errorf("branch %q not found", branch)
}
