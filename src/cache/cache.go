package cache

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
func (c *Cache[T]) Value() T { //nolint:ireturn
	if !c.initialized {
		panic("using a cached value before initialization")
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
