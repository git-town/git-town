package cache

import "github.com/git-town/git-town/v9/src/messages"

// Cache is a cache implementation for arbitrary data structures that ensures it is initialized.
// The zero value is an empty cache.
type Cache[T any] struct {
	initialized bool
	value       T
}

// Set allows collaborators to signal when the current branch has changed.
func (c *Cache[T]) Set(newValue T) {
	c.initialized = true
	c.value = newValue
}

// Value provides the current value.
func (c *Cache[T]) Value() T {
	if !c.initialized {
		panic(messages.CacheUsedUnitialized)
	}
	return c.value
}

// Initialized indicates if we have a current branch.
func (c *Cache[T]) Initialized() bool {
	return c.initialized
}

// Invalidate removes the cached value.
func (c *Cache[T]) Invalidate() {
	c.initialized = false
}
