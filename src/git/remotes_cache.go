package git

import (
	"errors"
)

// RemotesCache caches information about the current Git remotes.
type RemotesCache struct {
	remotes     []string
	initialized bool
}

// Get provides the currently cached remotes.
func (rc *RemotesCache) Get() ([]string, error) {
	if !rc.initialized {
		return rc.remotes, errors.New("cache not initialized")
	}
	return rc.remotes, nil
}

// Reset invalidates the cache.
func (rc *RemotesCache) Reset() {
	rc.initialized = false
}

// Set sets the cache to the given remotes.
func (rc *RemotesCache) Set(remotes []string) {
	rc.remotes = remotes
	rc.initialized = true
}
