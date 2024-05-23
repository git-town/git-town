package gitdomain

import (
	"fmt"

	. "github.com/git-town/git-town/v14/src/gohacks/prelude"
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

func EmptyBranchInfo() BranchInfo {
	return BranchInfo{
		LocalName:  None[LocalBranchName](),
		LocalSHA:   None[SHA](),
		RemoteName: None[RemoteBranchName](),
		RemoteSHA:  None[SHA](),
		SyncStatus: SyncStatusUpToDate,
	}
}

// TODO: rename to IsOmni after the existing IsOmni has been replaced with calls to this method
func (self BranchInfo) GetOmni() (bool, LocalBranchName, SHA) {
	localBranch, hasLocalBranch := self.LocalName.Get()
	localSHA, hasLocalSHA := self.LocalSHA.Get()
	remoteSHA, hasRemoteSHA := self.RemoteSHA.Get()
	isOmni := hasLocalBranch && hasLocalSHA && hasRemoteSHA && localSHA == remoteSHA
	return isOmni, localBranch, localSHA
}

// TODO: delete and replace with destructuring the LocalName property
func (self BranchInfo) HasLocalBranch() bool {
	return self.LocalName.IsSome() && self.LocalSHA.IsSome()
}

func (self BranchInfo) HasOnlyLocalBranch() bool {
	return self.HasLocalBranch() && !self.HasRemoteBranch()
}

func (self BranchInfo) HasOnlyRemoteBranch() bool {
	return self.HasRemoteBranch() && !self.HasLocalBranch()
}

// TODO: delete and replace with destructuring the RemoteName property
func (self BranchInfo) HasRemoteBranch() bool {
	return self.RemoteName.IsSome() && self.RemoteSHA.IsSome()
}

func (self BranchInfo) HasTrackingBranch() bool {
	return self.HasLocalBranch() && self.HasRemoteBranch()
}

// IsEmpty indicates whether this BranchInfo is completely empty, i.e. not a single branch contains something.
func (self BranchInfo) IsEmpty() bool {
	return !self.HasLocalBranch() && !self.HasRemoteBranch()
}

// TODO: delete and replace with destructuring the LocalName property
// IsLocalBranch indicates whether this branch exists in the local repo that Git Town is running in.
func (self BranchInfo) IsLocal() bool {
	return self.LocalName.IsSome()
}

// IsOmniBranch indicates whether the local and remote branch are in sync.
func (self BranchInfo) IsOmniBranch() bool {
	return !self.IsEmpty() && self.LocalSHA == self.RemoteSHA
}

func (self BranchInfo) String() string {
	return fmt.Sprintf("BranchInfo local: %s (%s) remote: %s (%s) %s", self.LocalName, self.LocalSHA, self.RemoteName, self.RemoteSHA, self.SyncStatus)
}
