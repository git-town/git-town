package git

// RemoteBranches caches the known remote branches of a Git repository.
// The zero value is valid.
type RemoteBranches struct {
	branches    []string
	initialized bool
}

// Initialized indicates whether this RemoteBranchTracker is initialized.
func (rbt *RemoteBranches) Initialized() bool {
	return rbt.initialized
}

// Get provides the currently known remote branches.
func (rbt *RemoteBranches) Get() []string {
	return rbt.branches
}

// Set stores the given remote branches.
func (rbt *RemoteBranches) Set(branches []string) {
	rbt.branches = branches
	rbt.initialized = true
}
