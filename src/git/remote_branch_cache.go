package git

// RemoteBranchCache caches the known remote branches of a Git repository.
// The zero value is valid.
type RemoteBranchCache struct {
	branches    []string
	initialized bool
}

// Initialized indicates whether this RemoteBranchTracker is initialized.
func (rbt *RemoteBranchCache) Initialized() bool {
	return rbt.initialized
}

// Get provides the currently known remote branches.
func (rbt *RemoteBranchCache) Get() []string {
	return rbt.branches
}

// Set stores the given remote branches.
func (rbt *RemoteBranchCache) Set(branches []string) {
	rbt.branches = branches
	rbt.initialized = true
}
