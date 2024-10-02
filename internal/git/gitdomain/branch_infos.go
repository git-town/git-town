package gitdomain

import (
	"fmt"

	"github.com/git-town/git-town/v16/internal/messages"
	. "github.com/git-town/git-town/v16/pkg/prelude"
)

// BranchInfos contains the BranchInfos for all branches in a repo.
// Tracking branches on the origin remote don't get their own entry,
// they are listed in the `TrackingBranch` property of the local branch they track.
type BranchInfos []BranchInfo

func (self BranchInfos) BranchIsActiveInAnotherWorktree(branch LocalBranchName) bool {
	branchInfo, has := self.FindByLocalName(branch).Get()
	if !has {
		return false
	}
	return branchInfo.SyncStatus == SyncStatusOtherWorktree
}

// FindByLocalName provides the branch with the given name if one exists.
func (self BranchInfos) FindByLocalName(branchName LocalBranchName) OptionP[BranchInfo] {
	for bi, branch := range self {
		if localName, hasLocalName := branch.LocalName.Get(); hasLocalName {
			if localName == branchName {
				return SomeP(&self[bi])
			}
		}
	}
	return NoneP[BranchInfo]()
}

// FindByRemoteName provides the local branch that has the given remote branch as its tracking branch
// or nil if no such branch exists.
func (self BranchInfos) FindByRemoteName(remoteBranch RemoteBranchName) OptionP[BranchInfo] {
	for b, bi := range self {
		if remoteName, hasRemoteName := bi.RemoteName.Get(); hasRemoteName {
			if remoteName == remoteBranch {
				return SomeP(&self[b])
			}
		}
	}
	return NoneP[BranchInfo]()
}

func (self BranchInfos) FindLocalOrRemote(branchName LocalBranchName) OptionP[BranchInfo] {
	branchInfoOpt := self.FindByLocalName(branchName)
	if branchInfoOpt.IsSome() {
		return branchInfoOpt
	}
	remoteName := branchName.AtRemote(RemoteOrigin)
	branchInfoOpt = self.FindByRemoteName(remoteName)
	if branchInfoOpt.IsSome() {
		return branchInfoOpt
	}
	return NoneP[BranchInfo]()
}

func (self BranchInfos) FindMatchingRecord(other BranchInfo) OptionP[BranchInfo] {
	for b, bi := range self {
		biLocalName, hasBiLocalName := bi.LocalName.Get()
		otherLocalName, hasOtherLocalName := other.LocalName.Get()
		if hasBiLocalName && hasOtherLocalName && biLocalName == otherLocalName {
			return SomeP(&self[b])
		}
		biRemoteName, hasBiRemoteName := bi.RemoteName.Get()
		otherRemoteName, hasOtherRemoteName := other.RemoteName.Get()
		if hasBiRemoteName && hasOtherRemoteName && biRemoteName == otherRemoteName {
			return SomeP(&self[b])
		}
	}
	return NoneP[BranchInfo]()
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
	return self.FindByRemoteName(localBranch.TrackingBranch()).IsSome()
}

// LocalBranches provides only the branches that exist on the local machine.
func (self BranchInfos) LocalBranches() BranchInfos {
	result := BranchInfos{}
	for _, bi := range self {
		if bi.LocalName.IsSome() {
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

// Select provides the BranchInfos with the given names.
func (self BranchInfos) Select(names ...LocalBranchName) (BranchInfos, error) {
	result := make(BranchInfos, len(names))
	for n, name := range names {
		if branchInfo, has := self.FindByLocalName(name).Get(); has {
			result[n] = *branchInfo
			continue
		}
		remoteName := name.AtRemote(RemoteOrigin)
		if branchInfo, has := self.FindByRemoteName(remoteName).Get(); has {
			result[n] = *branchInfo
			continue
		}
		return result, fmt.Errorf(messages.BranchDoesntExist, name)
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
