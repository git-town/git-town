package gitdomain

import (
	"fmt"

	"github.com/git-town/git-town/v12/src/messages"
)

// BranchInfos contains the BranchInfos for all branches in a repo.
// Tracking branches on the origin remote don't get their own entry,
// they are listed in the `TrackingBranch` property of the local branch they track.
type BranchInfos []BranchInfo

// FindByLocalName provides the branch with the given name if one exists.
func (self BranchInfos) FindByLocalName(branchName LocalBranchName) *BranchInfo {
	for bi, branch := range self {
		if branch.LocalName == branchName {
			return &self[bi]
		}
	}
	return nil
}

// FindByRemoteName provides the local branch that has the given remote branch as its tracking branch
// or nil if no such branch exists.
func (self BranchInfos) FindByRemoteName(remoteBranch RemoteBranchName) *BranchInfo {
	for b, bi := range self {
		if bi.RemoteName == remoteBranch {
			return &self[b]
		}
	}
	return nil
}

func (self BranchInfos) FindMatchingRecord(other BranchInfo) BranchInfo {
	for _, bi := range self {
		if bi.LocalName == other.LocalName && !other.LocalName.IsEmpty() {
			return bi
		}
		if bi.RemoteName == other.RemoteName && !other.RemoteName.IsEmpty() {
			return bi
		}
	}
	return EmptyBranchInfo()
}

// HasLocalBranch indicates whether the given local branch is already known to this BranchInfos instance.
func (self BranchInfos) HasLocalBranch(branch LocalBranchName) bool {
	for _, bi := range self {
		if bi.LocalName == branch {
			return true
		}
	}
	return false
}

// HasLocalBranches indicates whether this BranchInfos instance contains all the given branches.
func (self BranchInfos) HasLocalBranches(branches LocalBranchNames) bool {
	for _, branch := range branches {
		if !self.HasLocalBranch(branch) {
			return false
		}
	}
	return true
}

// HasMatchingRemoteBranchFor indicates whether there is already a remote branch matching the given local branch.
func (self BranchInfos) HasMatchingTrackingBranchFor(localBranch LocalBranchName) bool {
	return self.FindByRemoteName(localBranch.TrackingBranch()) != nil
}

// LocalBranches provides only the branches that exist on the local machine.
func (self BranchInfos) LocalBranches() BranchInfos {
	result := BranchInfos{}
	for _, bi := range self {
		if bi.IsLocal() {
			result = append(result, bi)
		}
	}
	return result
}

// LocalBranchesWithDeletedTrackingBranches provides only the branches that exist locally and have a deleted tracking branch.
func (self BranchInfos) LocalBranchesWithDeletedTrackingBranches() BranchInfos {
	result := BranchInfos{}
	for _, bi := range self {
		if bi.SyncStatus == SyncStatusDeletedAtRemote {
			result = append(result, bi)
		}
	}
	return result
}

// Names provides the names of all local branches in this BranchesSyncStatus instance.
func (self BranchInfos) Names() LocalBranchNames {
	result := make(LocalBranchNames, 0, len(self))
	for _, bi := range self {
		if !bi.LocalName.IsEmpty() {
			result = append(result, bi.LocalName)
		}
	}
	return result
}

func (self BranchInfos) Remove(branchName LocalBranchName) BranchInfos {
	result := BranchInfos{}
	for _, bi := range self {
		if bi.LocalName != branchName {
			result = append(result, bi)
		}
	}
	return result
}

// Select provides the BranchSyncStatus elements with the given names.
// TODO: make argument variadic
func (self BranchInfos) Select(names []LocalBranchName) (BranchInfos, error) {
	result := make(BranchInfos, len(names))
	for b, bi := range names {
		branch := self.FindByLocalName(bi)
		if branch == nil {
			return result, fmt.Errorf(messages.BranchDoesntExist, bi)
		}
		result[b] = *branch
	}
	return result, nil
}

func (self BranchInfos) UpdateLocalSHA(branch LocalBranchName, sha SHA) error {
	for b := range self {
		if self[b].LocalName == branch {
			self[b].LocalSHA = sha
			return nil
		}
	}
	return fmt.Errorf("branch %q not found", branch)
}
