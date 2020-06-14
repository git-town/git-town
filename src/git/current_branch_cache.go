package git

// CurrentBranchCache tracks the currently checked out branch in a git.Repo.
// The zero value is valid.
type CurrentBranchCache struct {
	value string
}

// Set allows collaborators to signal when the current branch has changed.
func (c *CurrentBranchCache) Set(newBranch string) {
	c.value = newBranch
}

// Current provides the currently checked out branch.
func (c *CurrentBranchCache) Current() string {
	if c.value == "" {
		panic("using current branch before initialization")
	}
	return c.value
}

// Initialized indicates if we have a current branch.
func (c *CurrentBranchCache) Initialized() bool {
	return c.value != ""
}

// Reset invalidates the cached value.
func (c *CurrentBranchCache) Reset() {
	c.value = ""
}
