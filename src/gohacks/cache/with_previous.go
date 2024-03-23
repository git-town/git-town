package cache

import "github.com/git-town/git-town/v13/src/messages"

// WithPrevious is a cache implementation for arbitrary data structures that keeps track of the current and previous values.
// The zero value is an empty cache.
type WithPrevious[T any] struct {
	initialized bool
	previous    T // the previous value
	value       T // the current value
}

// Initialized indicates if we have a current branch.
func (c *WithPrevious[T]) Initialized() bool {
	return c.initialized
}

// Invalidate removes the cached value.
func (c *WithPrevious[T]) Invalidate() {
	c.initialized = false
}

// Previous provides the previous value.
func (c *WithPrevious[T]) Previous() T {
	if !c.initialized {
		panic(messages.CacheUnitialized)
	}
	return c.previous
}

// Set allows collaborators to signal when the current branch has changed.
func (c *WithPrevious[T]) Set(newValue T) {
	c.previous = c.value
	c.value = newValue
	c.initialized = true
}

// Value provides the current value.
func (c *WithPrevious[T]) Value() T {
	if !c.initialized {
		panic(messages.CacheUnitialized)
	}
	return c.value
}
