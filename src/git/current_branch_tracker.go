package git

// CurrentBranchTracker tracks the currently checked out branch in a git.Repo.
// The zero value is valid.
type CurrentBranchTracker struct {
	value string
}

// Set allows collaborators to signal when the current branch has changed.
func (cbt *CurrentBranchTracker) Set(newBranch string) {
	cbt.value = newBranch
}

// Current provides the currently checked out branch.
func (cbt *CurrentBranchTracker) Current() string {
	if cbt.value == "" {
		panic("using current branch before initialization")
	}
	return cbt.value
}

// Initialized indicates if we have a current branch.
func (cbt *CurrentBranchTracker) Initialized() bool {
	return cbt.value != ""
}

// Reset invalidates the cached value.
func (cbt *CurrentBranchTracker) Reset() {
	cbt.value = ""
}
