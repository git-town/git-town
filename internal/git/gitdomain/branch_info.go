package gitdomain

import (
	"fmt"

	"github.com/git-town/git-town/v22/internal/messages"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

// BranchInfo describes the sync status of a branch in relation to its tracking branch.
type BranchInfo struct {
	Local Option[BranchData]

	// RemoteName contains the fully qualified name of the tracking branch, i.e. "origin/foo".
	RemoteName Option[RemoteBranchName]

	// RemoteSHA contains the SHA of the tracking branch before Git Town ran.
	RemoteSHA Option[SHA]

	// SyncStatus of the branch.
	SyncStatus SyncStatus
}

func (self BranchInfo) GetLocalOrRemoteName() BranchName {
	if local, hasLocal := self.Local.Get(); hasLocal {
		return local.Name.BranchName()
	}
	if remoteName, hasRemoteName := self.RemoteName.Get(); hasRemoteName {
		return remoteName.BranchName()
	}
	panic(messages.BranchInfoNoContent)
}

func (self BranchInfo) GetLocalOrRemoteNameAsLocalName() LocalBranchName {
	if local, hasLocal := self.Local.Get(); hasLocal {
		return local.Name
	}
	if remoteName, hasRemoteName := self.RemoteName.Get(); hasRemoteName {
		return remoteName.LocalBranchName()
	}
	panic(messages.BranchInfoNoContent)
}

func (self BranchInfo) GetLocalOrRemoteSHA() SHA {
	if local, has := self.Local.Get(); has {
		return local.SHA
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
	local, hasLocal := self.Local.Get()
	remote, hasRemote := self.RemoteSHA.Get()
	return hasLocal && hasRemote, local.SHA, remote
}

func (self BranchInfo) HasOnlyLocalBranch() bool {
	_, hasLocal := self.Local.Get()
	hasRemoteBranch, _, _ := self.GetRemote()
	return hasLocal && !hasRemoteBranch
}

func (self BranchInfo) HasOnlyRemoteBranch() bool {
	_, hasLocal := self.Local.Get()
	hasRemoteBranch, _, _ := self.GetRemote()
	return hasRemoteBranch && !hasLocal
}

func (self BranchInfo) HasTrackingBranch() bool {
	_, hasLocal := self.Local.Get()
	hasRemoteBranch, _, _ := self.GetRemote()
	return hasLocal && hasRemoteBranch
}

func (self BranchInfo) IsLocalOnlyBranch() (bool, LocalBranchName) {
	local, hasLocal := self.Local.Get()
	if !hasLocal {
		return false, ""
	}
	return self.RemoteName.IsNone(), local.Name
}

func (self BranchInfo) LocalName() Option[LocalBranchName] {
	if local, hasLocal := self.Local.Get(); hasLocal {
		return Some(local.Name)
	}
	return None[LocalBranchName]()
}

func (self BranchInfo) LocalSHA() Option[SHA] {
	if local, hasLocal := self.Local.Get(); hasLocal {
		return Some(local.SHA)
	}
	return None[SHA]()
}

// LocalBranchName provides the name of this branch as a local branch, independent of whether this branch is local or not.
func (self BranchInfo) LocalBranchName() LocalBranchName {
	if local, hasLocal := self.Local.Get(); hasLocal {
		return local.Name
	}
	if remoteName, hasRemoteName := self.RemoteName.Get(); hasRemoteName {
		return remoteName.LocalBranchName()
	}
	panic(messages.BranchInfoNoContent)
}

// IsOmniBranch indicates whether the branch described by this BranchInfo is omni
// and provides all relevant data around this scenario.
// An omni branch has the same SHA locally and remotely.
func (self BranchInfo) OmniBranch() Option[BranchData] {
	local, hasLocal := self.Local.Get()
	remoteSHA, hasRemoteSHA := self.RemoteSHA.Get()
	isOmni := hasLocal && hasRemoteSHA && local.SHA == remoteSHA
	if !isOmni {
		return None[BranchData]()
	}
	return Some(local)
}

func (self BranchInfo) String() string {
	local, hasLocal := self.Local.Get()
	if hasLocal {
		return fmt.Sprintf("BranchInfo local: %s (%s) remote: %s (%s) %s", local.Name, local.SHA, self.RemoteName, self.RemoteSHA, self.SyncStatus)
	}
	return fmt.Sprintf("BranchInfo local: (none) remote: %s (%s) %s", self.RemoteName, self.RemoteSHA, self.SyncStatus)
}
