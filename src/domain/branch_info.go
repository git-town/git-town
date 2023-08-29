package domain

import (
	"fmt"
)

// BranchInfo describes the sync status of a branch in relation to its tracking branch.
type BranchInfo struct {
	// Name contains the local name of the branch.
	Name LocalBranchName

	// InitialSHA contains the SHA that this branch had locally before Git Town ran.
	// TODO: rename to LocalSHA
	InitialSHA SHA

	// SyncStatus of the branch
	SyncStatus SyncStatus

	// RemoteName contains the fully qualified name of the tracking branch, i.e. "origin/foo".
	RemoteName RemoteBranchName

	// RemoteSHA contains the SHA of the tracking branch before Git Town ran.
	RemoteSHA SHA
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

// IsLocalBranch indicates whether this branch exists in the local repo that Git Town is running in.
func (bi BranchInfo) IsLocal() bool {
	return bi.SyncStatus.IsLocal()
}
