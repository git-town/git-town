package git

type SyncStatus int

const (
	SyncStatusUpToDate SyncStatus = iota
	SyncStatusBehind
	SyncStatusAhead
)
