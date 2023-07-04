package git

// BranchLocation encodes the various places that a branch can exist at.
type BranchLocation int

const (
	BranchLocationLocalAndRemote  BranchLocation = iota // the branch exists locally and at the remote
	BranchLocationLocalOnly                             // the branch was created locally and hasn't been pushed to the remote yet
	BranchLocationRemoteOnly                            // the branch exists only at the remote
	BranchLocationDeletedAtRemote                       // the branch was deleted on the remote
)
