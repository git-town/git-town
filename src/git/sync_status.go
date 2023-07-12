package git

// SyncStatus encodes the places a branch can exist at.
type SyncStatus int

const (
	SyncStatusUpToDate        SyncStatus = iota // the branch exists locally and remotely, the local branch is up to date
	SyncStatusBehind                            // the branch exists locally and remotely, the local branch is behind the remote tracking branch
	SyncStatusAhead                             // the branch exists locally and remotely, the local branch is ahead of its remote branch
	SyncStatusLocalOnly                         // the branch was created locally and hasn't been pushed to the remote yet
	SyncStatusRemoteOnly                        // the branch exists only at the remote
	SyncStatusDeletedAtRemote                   // the branch was deleted on the remote
)
