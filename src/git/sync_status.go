package git

import (
	"fmt"
)

// SyncStatus encodes the places a branch can exist at.
// This is a type-safe enum, see https://npf.io/2022/05/safer-enums.
type SyncStatus struct {
	name string
}

func (s SyncStatus) String() string { return s.name }

var (
	SyncStatusUpToDate        = SyncStatus{"up to date"}        // the branch exists locally and remotely, the local branch is up to date
	SyncStatusBehind          = SyncStatus{"behind"}            // the branch exists locally and remotely, the local branch is behind the remote tracking branch
	SyncStatusAhead           = SyncStatus{"ahead"}             // the branch exists locally and remotely, the local branch is ahead of its remote branch
	SyncStatusLocalOnly       = SyncStatus{"local only"}        // the branch was created locally and hasn't been pushed to the remote yet
	SyncStatusRemoteOnly      = SyncStatus{"remote only"}       // the branch exists only at the remote
	SyncStatusDeletedAtRemote = SyncStatus{"deleted at remote"} // the branch was deleted on the remote
)

// IsLocal indicates whether a branch with this SyncStatus exists in the local repo.
func (s SyncStatus) IsLocal() bool {
	switch s {
	case SyncStatusLocalOnly, SyncStatusUpToDate, SyncStatusAhead, SyncStatusBehind, SyncStatusDeletedAtRemote:
		return true
	case SyncStatusRemoteOnly:
		return false
	}
	panic(fmt.Sprintf("uncaptured sync status: %v", s))
}
