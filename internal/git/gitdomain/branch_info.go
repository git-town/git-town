package gitdomain

import (
	"fmt"

	. "github.com/git-town/git-town/v14/internal/gohacks/prelude"
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

	// SyncStatus of the branch
	SyncStatus SyncStatus
}

// provides both the name and SHA of the local branch
func (self BranchInfo) GetLocal() (bool, LocalBranchName, SHA) {
	name, hasName := self.LocalName.Get()
	sha, hasSHA := self.LocalSHA.Get()
	return hasName && hasSHA, name, sha
}

// provides both the name and SHA of the remote branch
func (self BranchInfo) GetRemote() (bool, RemoteBranchName, SHA) {
	name, hasName := self.RemoteName.Get()
	sha, hasSHA := self.RemoteSHA.Get()
	return hasName && hasSHA, name, sha
}

// provides the SHAs of the local and remote branch
func (self BranchInfo) GetSHAs() (hasBothSHA bool, localSHA, remoteSHA SHA) {
	local, hasLocal := self.LocalSHA.Get()
	remote, hasRemote := self.RemoteSHA.Get()
	return hasLocal && hasRemote, local, remote
}

func (self BranchInfo) HasLocalBranch() (hasLocalBranch bool, branchName LocalBranchName, sha SHA) {
	localName, hasLocalName := self.LocalName.Get()
	localSHA, hasLocalSHA := self.LocalSHA.Get()
	hasLocalBranch = hasLocalName && hasLocalSHA
	return hasLocalBranch, localName, localSHA
}

func (self BranchInfo) HasOnlyLocalBranch() bool {
	hasLocalBranch, _, _ := self.HasLocalBranch()
	hasRemoteBranch, _, _ := self.HasRemoteBranch()
	return hasLocalBranch && !hasRemoteBranch
}

func (self BranchInfo) HasOnlyRemoteBranch() bool {
	hasLocalBranch, _, _ := self.HasLocalBranch()
	hasRemoteBranch, _, _ := self.HasRemoteBranch()
	return hasRemoteBranch && !hasLocalBranch
}

func (self BranchInfo) HasRemoteBranch() (hasRemoteBranch bool, remoteBranchName RemoteBranchName, remoteBranchSHA SHA) {
	remoteName, hasRemoteName := self.RemoteName.Get()
	remoteSHA, hasRemoteSHA := self.RemoteSHA.Get()
	hasRemoteBranch = hasRemoteName && hasRemoteSHA
	return hasRemoteBranch, remoteName, remoteSHA
}

func (self BranchInfo) HasTrackingBranch() bool {
	hasLocalBranch, _, _ := self.HasLocalBranch()
	hasRemoteBranch, _, _ := self.HasRemoteBranch()
	return hasLocalBranch && hasRemoteBranch
}

// Indicates whether the branch described by this BranchInfo is omni
// and provides all relevant data around this scenario.
// An omni branch has the same SHA locally and remotely.
func (self BranchInfo) IsOmniBranch() (isOmni bool, branch LocalBranchName, sha SHA) {
	localSHA, hasLocalSHA := self.LocalSHA.Get()
	branchName, hasBranch := self.LocalName.Get()
	remoteSHA, hasRemoteSHA := self.RemoteSHA.Get()
	isOmni = hasLocalSHA && hasRemoteSHA && hasBranch && localSHA == remoteSHA
	return isOmni, branchName, localSHA
}

func (self BranchInfo) String() string {
	return fmt.Sprintf("BranchInfo local: %s (%s) remote: %s (%s) %s", self.LocalName, self.LocalSHA, self.RemoteName, self.RemoteSHA, self.SyncStatus)
}
