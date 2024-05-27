package gitdomain

import (
	"fmt"

	. "github.com/git-town/git-town/v14/src/gohacks/prelude"
	"github.com/git-town/git-town/v14/src/messages"
)

// BranchInfos contains the BranchInfos for all branches in a repo.
// Tracking branches on the origin remote don't get their own entry,
// they are listed in the `TrackingBranch` property of the local branch they track.
type BranchInfos []BranchInfo

// FindByLocalName provides the branch with the given name if one exists.
func (self BranchInfos) FindByLocalName(branchName LocalBranchName) Option[BranchInfo] {
	for bi, branch := range self {
		if localName, hasLocalName := branch.LocalName.Get(); hasLocalName {
			if localName == branchName {
				return Some(self[bi])
			}
		}
	}
	return None[BranchInfo]()
}

// FindByRemoteName provides the local branch that has the given remote branch as its tracking branch
// or nil if no such branch exists.
func (self BranchInfos) FindByRemoteName(remoteBranch RemoteBranchName) *BranchInfo {
	for b, bi := range self {
		if remoteName, hasRemoteName := bi.RemoteName.Get(); hasRemoteName {
			if remoteName == remoteBranch {
				return &self[b]
			}
		}
	}
	return nil
}

func (self BranchInfos) FindMatchingRecord(other BranchInfo) BranchInfo {
	for _, bi := range self {
		biLocalName, hasBiLocalName := bi.LocalName.Get()
		otherLocalName, hasOtherLocalName := other.LocalName.Get()
		if hasBiLocalName && hasOtherLocalName && biLocalName == otherLocalName {
			return bi
		}
		biRemoteName, hasBiRemoteName := bi.RemoteName.Get()
		otherRemoteName, hasOtherRemoteName := other.RemoteName.Get()
		if hasBiRemoteName && hasOtherRemoteName && biRemoteName == otherRemoteName {
			return bi
		}
	}
	return EmptyBranchInfo()
}

// HasLocalBranch indicates whether the given local branch is already known to this BranchInfos instance.
func (self BranchInfos) HasLocalBranch(branch LocalBranchName) bool {
	for _, bi := range self {
		if biLocalName, hasBiLocalName := bi.LocalName.Get(); hasBiLocalName {
			if biLocalName == branch {
				return true
			}
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
		if localName, hasLocalName := bi.LocalName.Get(); hasLocalName {
			result = append(result, localName)
		}
	}
	return result
}

func (self BranchInfos) Remove(branchName LocalBranchName) BranchInfos {
	result := BranchInfos{}
	for _, bi := range self {
		localName, hasLocalName := bi.LocalName.Get()
		if !hasLocalName || localName != branchName {
			result = append(result, bi)
		}
	}
	return result
}

// Select provides the BranchSyncStatus elements with the given names.
func (self BranchInfos) Select(names ...LocalBranchName) (BranchInfos, error) {
	result := make(BranchInfos, len(names))
	for b, bi := range names {
		if branch, hasBranch := self.FindByLocalName(bi).Get(); hasBranch {
			result[b] = branch
		} else {
			return result, fmt.Errorf(messages.BranchDoesntExist, bi)
		}
	}
	return result, nil
}

func (self BranchInfos) UpdateLocalSHA(branch LocalBranchName, sha SHA) error {
	for b := range self {
		if localName, hasLocalName := self[b].LocalName.Get(); hasLocalName {
			if localName == branch {
				self[b].LocalSHA = Some(sha)
				return nil
			}
		}
	}
	return fmt.Errorf("branch %q not found", branch)
}
