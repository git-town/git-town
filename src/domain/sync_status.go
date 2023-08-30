package domain

import (
	"encoding/json"
	"fmt"
)

// SyncStatus encodes the places a branch can exist at.
// This is a type-safe enum, see https://npf.io/2022/05/safer-enums.
type SyncStatus struct {
	name string
}

func (s SyncStatus) String() string { return s.name }

var (
	SyncStatusUpToDate        = SyncStatus{"up to date"}        //nolint:gochecknoglobals // the branch exists locally and remotely, the local branch is up to date
	SyncStatusBehind          = SyncStatus{"behind"}            //nolint:gochecknoglobals // the branch exists locally and remotely, the local branch is behind the remote tracking branch
	SyncStatusAhead           = SyncStatus{"ahead"}             //nolint:gochecknoglobals // the branch exists locally and remotely, the local branch is ahead of its remote branch
	SyncStatusAheadAndBehind  = SyncStatus{"ahead and behind"}  //nolint:gochecknoglobals // the branch exists locally and remotely, both ends have different commits
	SyncStatusLocalOnly       = SyncStatus{"local only"}        //nolint:gochecknoglobals // the branch was created locally and hasn't been pushed to the remote yet
	SyncStatusRemoteOnly      = SyncStatus{"remote only"}       //nolint:gochecknoglobals // the branch exists only at the remote
	SyncStatusDeletedAtRemote = SyncStatus{"deleted at remote"} //nolint:gochecknoglobals // the branch was deleted on the remote
)

// IsLocal indicates whether a branch with this SyncStatus exists in the local repo.
func (s SyncStatus) IsLocal() bool {
	switch s {
	case SyncStatusLocalOnly, SyncStatusUpToDate, SyncStatusAhead, SyncStatusBehind, SyncStatusAheadAndBehind, SyncStatusDeletedAtRemote:
		return true
	case SyncStatusRemoteOnly:
		return false
	}
	panic(fmt.Sprintf("uncaptured sync status: %v", s))
}

// MarshalJSON is used when serializing this LocalBranchName to JSON.
func (s SyncStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.name)
}

// UnmarshalJSON is used when de-serializing JSON into a LocalBranchName.
func (s *SyncStatus) UnmarshalJSON(ba []byte) error {
	return json.Unmarshal(ba, &s.name)
}
