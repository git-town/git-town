package gitdomain

import (
	"fmt"

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
func (self BranchInfos) FindByLocalName(branchName LocalBranchName) OptionalMutable[BranchInfo] {
	for bi, branch := range self {
		if localName, hasLocalName := branch.LocalName.Get(); hasLocalName {
			if localName == branchName {
				return MutableSome(&self[bi])
			}
		}
	}
	return MutableNone[BranchInfo]()
}

// FindByRemoteName provides the local branch that has the given remote branch as its tracking branch
// or nil if no such branch exists.
func (self BranchInfos) FindByRemoteName(remoteBranch RemoteBranchName) OptionalMutable[BranchInfo] {
	for b, bi := range self {
		if remoteName, hasRemoteName := bi.RemoteName.Get(); hasRemoteName {
			if remoteName == remoteBranch {
				return MutableSome(&self[b])
			}
		}
	}
	return MutableNone[BranchInfo]()
}

func (self BranchInfos) FindLocalOrRemote(branchName LocalBranchName, remote Remote) OptionalMutable[BranchInfo] {
	branchInfoOpt := self.FindByLocalName(branchName)
	if branchInfoOpt.IsSome() {
		return branchInfoOpt
	}
	remoteName := branchName.AtRemote(remote)
	branchInfoOpt = self.FindByRemoteName(remoteName)
	if branchInfoOpt.IsSome() {
		return branchInfoOpt
	}
	return MutableNone[BranchInfo]()
}

func (self BranchInfos) FindMatchingRecord(other BranchInfo) OptionalMutable[BranchInfo] {
	for b, bi := range self {
		biLocalName, hasBiLocalName := bi.LocalName.Get()
		otherLocalName, hasOtherLocalName := other.LocalName.Get()
		if hasBiLocalName && hasOtherLocalName && biLocalName == otherLocalName {
			return MutableSome(&self[b])
		}
		biRemoteName, hasBiRemoteName := bi.RemoteName.Get()
		otherRemoteName, hasOtherRemoteName := other.RemoteName.Get()
		if hasBiRemoteName && hasOtherRemoteName && biRemoteName == otherRemoteName {
			return MutableSome(&self[b])
		}
	}
	return MutableNone[BranchInfo]()
}

func (self BranchInfos) HasBranch(branch LocalBranchName) bool {
	for _, branchInfo := range self {
		if localName, hasLocalName := branchInfo.LocalName.Get(); hasLocalName {
			if localName == branch {
				return true
			}
		}
		if trackingName, hasTrackingBranch := branchInfo.RemoteName.Get(); hasTrackingBranch {
			if trackingName.LocalBranchName() == branch {
				return true
			}
		}
	}
	return false
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
func (self BranchInfos) HasMatchingTrackingBranchFor(localBranch LocalBranchName, devRemote Remote) bool {
	return self.FindByRemoteName(localBranch.TrackingBranch(devRemote)).IsSome()
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
func (self BranchInfos) Select(remote Remote, names ...LocalBranchName) (result BranchInfos, nonExisting LocalBranchNames) {
	result = make(BranchInfos, 0, len(names))
	for _, name := range names {
		if branchInfo, has := self.FindByLocalName(name).Get(); has {
			result = append(result, *branchInfo)
			continue
		}
		remoteName := name.AtRemote(remote)
		if branchInfo, has := self.FindByRemoteName(remoteName).Get(); has {
			result = append(result, *branchInfo)
			continue
		}
		nonExisting = append(nonExisting, name)
	}
	return result, nonExisting
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
