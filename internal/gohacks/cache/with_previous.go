package cache

import "github.com/git-town/git-town/v21/internal/messages"

// WithPrevious is a cache implementation for arbitrary data structures that keeps track of the current and previous values.
// The zero value is an empty cache.
type WithPrevious[T any] struct {
	initialized bool
	previous    T // the previous value
	value       T // the current value
}

// Initialized indicates if we have a current branch.
func (self *WithPrevious[T]) Initialized() bool {
	return self.initialized
}

// Invalidate removes the cached value.
func (self *WithPrevious[T]) Invalidate() {
	self.initialized = false
}

// Previous provides the previous value.
func (self *WithPrevious[T]) Previous() T {
	if !self.initialized {
		panic(messages.CacheUnitialized)
	}
	return self.previous
}

// Set allows collaborators to signal when the current branch has changed.
func (self *WithPrevious[T]) Set(newValue T) {
	self.previous = self.value
	self.value = newValue
	self.initialized = true
}

// Value provides the current value.
func (self *WithPrevious[T]) Value() T {
	if !self.initialized {
		panic(messages.CacheUnitialized)
	}
	return self.value
}
