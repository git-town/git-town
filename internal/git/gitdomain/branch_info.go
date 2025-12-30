package gitdomain

import (
	"fmt"

	"github.com/git-town/git-town/v22/internal/messages"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

// BranchInfo describes the sync status of a branch in relation to its tracking branch.
type BranchInfo struct {
	// LocalName contains the local name of the branch.
	LocalName Option[LocalBranchName]

	// LocalSHA contains the SHA that this branch had locally before Git Town ran.
	LocalSHA Option[SHA]

	// RemoteName contains the fully qualified name of the tracking branch, i.e. "origin/foo".
	RemoteName Option[RemoteBranchName]

	// RemoteSHA contains the SHA of the tracking branch before Git Town ran.
	RemoteSHA Option[SHA]

	// SyncStatus of the branch.
	SyncStatus SyncStatus
}

// GetLocal provides both the name and SHA of the local branch.
func (self BranchInfo) GetLocal() (bool, LocalBranchName, SHA) {
	name, hasName := self.LocalName.Get()
	sha, hasSHA := self.LocalSHA.Get()
	return hasName && hasSHA, name, sha
}

func (self BranchInfo) GetLocalOrRemoteName() BranchName {
	if localName, hasLocalName := self.LocalName.Get(); hasLocalName {
		return localName.BranchName()
	}
	if remoteName, hasRemoteName := self.RemoteName.Get(); hasRemoteName {
		return remoteName.BranchName()
	}
	panic(messages.BranchInfoNoContent)
}

func (self BranchInfo) GetLocalOrRemoteNameAsLocalName() LocalBranchName {
	if localName, hasLocalName := self.LocalName.Get(); hasLocalName {
		return localName
	}
	if remoteName, hasRemoteName := self.RemoteName.Get(); hasRemoteName {
		return remoteName.LocalBranchName()
	}
	panic(messages.BranchInfoNoContent)
}

func (self BranchInfo) GetLocalOrRemoteSHA() SHA {
	if localSHA, has := self.LocalSHA.Get(); has {
		return localSHA
	}
	if remoteSHA, has := self.RemoteSHA.Get(); has {
		return remoteSHA
	}
	panic(messages.BranchInfoNoContent)
}

// GetRemote provides both the name and SHA of the remote branch.
func (self BranchInfo) GetRemote() (bool, RemoteBranchName, SHA) {
	name, hasName := self.RemoteName.Get()
	sha, hasSHA := self.RemoteSHA.Get()
	return hasName && hasSHA, name, sha
}

// GetSHAs provides the SHAs of the local and remote branch.
func (self BranchInfo) GetSHAs() (hasBothSHA bool, localSHA, remoteSHA SHA) {
	local, hasLocal := self.LocalSHA.Get()
	remote, hasRemote := self.RemoteSHA.Get()
	return hasLocal && hasRemote, local, remote
}

func (self BranchInfo) HasOnlyLocalBranch() bool {
	hasLocalBranch, _, _ := self.GetLocal()
	hasRemoteBranch, _, _ := self.GetRemote()
	return hasLocalBranch && !hasRemoteBranch
}

func (self BranchInfo) HasOnlyRemoteBranch() bool {
	hasLocalBranch, _, _ := self.GetLocal()
	hasRemoteBranch, _, _ := self.GetRemote()
	return hasRemoteBranch && !hasLocalBranch
}

func (self BranchInfo) HasTrackingBranch() bool {
	hasLocalBranch, _, _ := self.GetLocal()
	hasRemoteBranch, _, _ := self.GetRemote()
	return hasLocalBranch && hasRemoteBranch
}

func (self BranchInfo) IsLocalOnlyBranch() (bool, LocalBranchName) {
	branchName, hasLocalBranch := self.LocalName.Get()
	if !hasLocalBranch {
		return false, branchName
	}
	return self.RemoteName.IsNone(), branchName
}

// IsOmniBranch indicates whether the branch described by this BranchInfo is omni
// and provides all relevant data around this scenario.
// An omni branch has the same SHA locally and remotely.
func (self BranchInfo) IsOmniBranch() (isOmni bool, branch LocalBranchName, sha SHA) {
	localSHA, hasLocalSHA := self.LocalSHA.Get()
	branchName, hasBranch := self.LocalName.Get()
	remoteSHA, hasRemoteSHA := self.RemoteSHA.Get()
	isOmni = hasLocalSHA && hasRemoteSHA && hasBranch && localSHA == remoteSHA
	return isOmni, branchName, localSHA
}

// LocalBranchName provides the name of this branch as a local branch, independent of whether this branch is local or not.
func (self BranchInfo) LocalBranchName() LocalBranchName {
	if localName, hasLocalName := self.LocalName.Get(); hasLocalName {
		return localName
	}
	if remoteName, hasRemoteName := self.RemoteName.Get(); hasRemoteName {
		return remoteName.LocalBranchName()
	}
	panic(messages.BranchInfoNoContent)
}

func (self BranchInfo) String() string {
	return fmt.Sprintf("BranchInfo local: %s (%s) remote: %s (%s) %s", self.LocalName, self.LocalSHA, self.RemoteName, self.RemoteSHA, self.SyncStatus)
}
