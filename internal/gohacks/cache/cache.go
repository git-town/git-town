package cache

import "github.com/git-town/git-town/v20/internal/messages"

// Cache is a cache implementation for arbitrary data structures that ensures it is initialized.
// The zero value is an empty cache.
type Cache[T any] struct {
	initialized bool
	value       *T
}

// Initialized indicates if we have a current branch.
func (self *Cache[T]) Initialized() bool {
	return self.initialized
}

// Invalidate removes the cached value.
func (self *Cache[T]) Invalidate() {
	self.initialized = false
}

// Set allows collaborators to signal when the current branch has changed.
func (self *Cache[T]) Set(newValue *T) {
	self.value = newValue
	self.initialized = true
}

// Value provides the current value.
func (self *Cache[T]) Value() *T {
	if !self.initialized {
		panic(messages.CacheUnitialized)
	}
	return self.value
}
