package git

// RemoteBranchTracker caches the known remote branches of a Git repository.
// The zero value is valid.
type RemoteBranchTracker struct {
	branches    []string
	initialized bool
}

// Initialized indicates whether this RemoteBranchTracker is initialized.
func (rbt *RemoteBranchTracker) Initialized() bool {
	return rbt.initialized
}

// Branches provides the currently known remote branches.
func (rbt *RemoteBranchTracker) Branches() []string {
	return rbt.branches
}

// Set stores the given remote branches.
func (rbt *RemoteBranchTracker) Set(branches []string) {
	rbt.branches = branches
	rbt.initialized = true
}
