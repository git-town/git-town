package git

type rootDirectoryCache struct {
	value string
}

// Set allows collaborators to signal when the current branch has changed.
func (r *rootDirectoryCache) Set(newDir string) {
	r.value = newDir
}

// Current provides the currently checked out branch.
func (r *rootDirectoryCache) Current() string {
	if r.value == "" {
		panic("using current branch before initialization")
	}
	return r.value
}

// Initialized indicates if we have a current branch.
func (r *rootDirectoryCache) Initialized() bool {
	return r.value != ""
}

// Reset invalidates the cached value.
func (r *rootDirectoryCache) Reset() {
	r.value = ""
}
