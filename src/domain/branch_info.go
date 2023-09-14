package domain

import (
	"fmt"
)

// BranchInfo describes the sync status of a branch in relation to its tracking branch.
type BranchInfo struct {
	// LocalName contains the local name of the branch.
	LocalName LocalBranchName

	// LocalSHA contains the SHA that this branch had locally before Git Town ran.
	LocalSHA SHA

	// SyncStatus of the branch
	SyncStatus SyncStatus

	// RemoteName contains the fully qualified name of the tracking branch, i.e. "origin/foo".
	RemoteName RemoteBranchName

	// RemoteSHA contains the SHA of the tracking branch before Git Town ran.
	RemoteSHA SHA
}

func EmptyBranchInfo() BranchInfo {
	return BranchInfo{
		LocalName:  LocalBranchName{},
		LocalSHA:   EmptySHA(),
		SyncStatus: SyncStatusUpToDate,
		RemoteName: EmptyRemoteBranchName(),
		RemoteSHA:  EmptySHA(),
	}
}

func (bi BranchInfo) HasLocalBranch() bool {
	return !bi.LocalName.IsEmpty() && !bi.LocalSHA.IsEmpty()
}

func (bi BranchInfo) HasOnlyLocalBranch() bool {
	return bi.HasLocalBranch() && !bi.HasRemoteBranch()
}

func (bi BranchInfo) HasOnlyRemoteBranch() bool {
	return bi.HasRemoteBranch() && !bi.HasLocalBranch()
}

func (bi BranchInfo) HasRemoteBranch() bool {
	return !bi.RemoteName.IsEmpty() && !bi.RemoteSHA.IsEmpty()
}

func (bi BranchInfo) HasTrackingBranch() bool {
	switch bi.SyncStatus {
	case SyncStatusAhead, SyncStatusBehind, SyncStatusAheadAndBehind, SyncStatusUpToDate, SyncStatusRemoteOnly:
		return true
	case SyncStatusLocalOnly, SyncStatusDeletedAtRemote:
		return false
	}
	panic(fmt.Sprintf("unknown sync status: %v", bi.SyncStatus))
}

func (bi BranchInfo) IsEmpty() bool {
	return bi.LocalName.IsEmpty() && bi.LocalSHA.IsEmpty() && bi.RemoteName.IsEmpty() && bi.RemoteSHA.IsEmpty()
}

// IsLocalBranch indicates whether this branch exists in the local repo that Git Town is running in.
func (bi BranchInfo) IsLocal() bool {
	return bi.SyncStatus.IsLocal()
}

// IsOmniBranch indicates whether the local and remote branch are in sync.
func (bi BranchInfo) IsOmniBranch() bool {
	return !bi.IsEmpty() && bi.LocalSHA == bi.RemoteSHA
}
