package git

// CurrentBranchCache tracks the currently checked out branch in a git.Repo.
// The zero value is valid.
type CurrentBranchCache struct {
	value string
}

// Set allows collaborators to signal when the current branch has changed.
func (cbt *CurrentBranchCache) Set(newBranch string) {
	cbt.value = newBranch
}

// Current provides the currently checked out branch.
func (cbt *CurrentBranchCache) Current() string {
	if cbt.value == "" {
		panic("using current branch before initialization")
	}
	return cbt.value
}

// Initialized indicates if we have a current branch.
func (cbt *CurrentBranchCache) Initialized() bool {
	return cbt.value != ""
}

// Reset invalidates the cached value.
func (cbt *CurrentBranchCache) Reset() {
	cbt.value = ""
}
