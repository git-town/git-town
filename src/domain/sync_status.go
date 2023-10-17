package domain

import (
	"fmt"
)

// SyncStatus encodes the places a branch can exist at.
// This is a type-safe enum, see https://npf.io/2022/05/safer-enums.
type SyncStatus struct {
	name string
}

// IsLocal indicates whether a branch with this SyncStatus exists in the local repo.
func (self SyncStatus) IsLocal() bool {
	switch self {
	case SyncStatusLocalOnly, SyncStatusUpToDate, SyncStatusAhead, SyncStatusBehind, SyncStatusAheadAndBehind, SyncStatusDeletedAtRemote:
		return true
	case SyncStatusRemoteOnly:
		return false
	}
	panic(fmt.Sprintf("uncaptured sync status: %v", self))
}

func (self SyncStatus) String() string { return self.name }

var (
	SyncStatusUpToDate        = SyncStatus{"up to date"}        //nolint:gochecknoglobals // the branch exists locally and remotely, the local branch is up to date
	SyncStatusBehind          = SyncStatus{"behind"}            //nolint:gochecknoglobals // the branch exists locally and remotely, the local branch is behind the remote tracking branch
	SyncStatusAhead           = SyncStatus{"ahead"}             //nolint:gochecknoglobals // the branch exists locally and remotely, the local branch is ahead of its remote branch
	SyncStatusAheadAndBehind  = SyncStatus{"ahead and behind"}  //nolint:gochecknoglobals // the branch exists locally and remotely, both ends have different commits
	SyncStatusLocalOnly       = SyncStatus{"local only"}        //nolint:gochecknoglobals // the branch was created locally and hasn't been pushed to the remote yet
	SyncStatusRemoteOnly      = SyncStatus{"remote only"}       //nolint:gochecknoglobals // the branch exists only at the remote
	SyncStatusDeletedAtRemote = SyncStatus{"deleted at remote"} //nolint:gochecknoglobals // the branch was deleted on the remote
)
