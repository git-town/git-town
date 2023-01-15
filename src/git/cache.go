package git

// Cache stores data of the generic type.
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
		panic("using current branch before initialization")
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
