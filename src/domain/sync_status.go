package domain

// SyncStatus encodes the places a branch can exist at.
// This is a type-safe enum, see https://npf.io/2022/05/safer-enums.
type SyncStatus string

func (self SyncStatus) String() string {
	return string(self)
}

const (
	SyncStatusUpToDate        SyncStatus = "up to date"                 // the branch exists locally and remotely, the local branch is up to date
	SyncStatusNotInSync       SyncStatus = "not in sync"                // the branch exists locally and remotely, the local branch is behind the remote tracking branch
	SyncStatusLocalOnly       SyncStatus = "local only"                 // the branch was created locally and hasn't been pushed to the remote yet
	SyncStatusRemoteOnly      SyncStatus = "remote only"                // the branch exists only at the remote
	SyncStatusDeletedAtRemote SyncStatus = "deleted at remote"          // the branch was deleted on the remote
	SyncStatusOtherWorktree   SyncStatus = "active in another worktree" // the branch is active in another worktree and should not be synced
)
