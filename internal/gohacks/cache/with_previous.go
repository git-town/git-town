package cache

import . "github.com/git-town/git-town/v22/pkg/prelude"

// WithPrevious is a cache implementation for arbitrary data structures that keeps track of the current and previous values.
// The zero value is an empty cache.
type WithPrevious[T any] struct {
	previous Option[T] // the previous value
	value    Option[T] // the current value
}

// Get provides the current value.
func (self *WithPrevious[T]) Get() (T, bool) {
	return self.value.Get()
}

// GetPrevious provides the previous value.
func (self *WithPrevious[T]) GetPrevious() (T, bool) {
	return self.previous.Get()
}

// Invalidate removes all cached values.
func (self *WithPrevious[T]) Invalidate() {
	self.previous = None[T]()
	self.value = None[T]()
}

// Set assigns a new current value.
func (self *WithPrevious[T]) Set(newValue T) {
	self.previous = self.value
	self.value = Some(newValue)
}
