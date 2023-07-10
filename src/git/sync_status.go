package git

// SyncStatus describes whether a Git branch is up to date with its tracking branch.
type SyncStatus int8

const (
	SyncStatusUpToDate SyncStatus = iota
	SyncStatusBehind
	SyncStatusAhead
)
