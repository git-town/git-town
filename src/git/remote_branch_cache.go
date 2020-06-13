package git

import "errors"

// RemoteBranchCache caches the known remote branches of a Git repository.
// The zero value is valid.
type RemoteBranchCache struct {
	branches    []string
	initialized bool
}

// Get provides the currently known remote branches.
func (rbt *RemoteBranchCache) Get() ([]string, error) {
	if !rbt.initialized {
		return rbt.branches, errors.New("not initialized")
	}
	return rbt.branches, nil
}

// Initialized stores the given remote branches.
func (rbt *RemoteBranchCache) Initialized() bool {
	return rbt.initialized
}

// Reset stores the given remote branches.
func (rbt *RemoteBranchCache) Reset(branches []string) {
	rbt.initialized = false
}

// Set stores the given remote branches.
func (rbt *RemoteBranchCache) Set(branches []string) {
	rbt.branches = branches
	rbt.initialized = true
}
