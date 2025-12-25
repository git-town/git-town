package cache

import . "github.com/git-town/git-town/v22/pkg/prelude"

// Cache is a cache implementation for arbitrary data structures that ensures it is initialized.
// The zero value is an empty cache.
type Cache[T any] struct {
	value Option[T]
}

// Get provides the current value.
func (self *Cache[T]) Get() (T, bool) {
	return self.value.Get()
}

// Invalidate removes the cached value.
func (self *Cache[T]) Invalidate() {
	self.value = None[T]()
}

// Set allows collaborators to signal when the current branch has changed.
func (self *Cache[T]) Set(newValue T) {
	self.value = Some(newValue)
}
