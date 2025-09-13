package cache

import . "github.com/git-town/git-town/v21/pkg/prelude"

// WithPrevious is a cache implementation for arbitrary data structures that keeps track of the current and previous values.
// The zero value is an empty cache.
type WithPrevious[T any] struct {
	previous Option[T] // the previous value
	value    Option[T] // the current value
}

// Invalidate removes all cached values.
func (self *WithPrevious[T]) Invalidate() {
	self.previous = None[T]()
	self.value = None[T]()
}

// Previous provides the previous value.
func (self *WithPrevious[T]) Previous() Option[T] {
	return self.previous
}

// Sets a new current value.
func (self *WithPrevious[T]) Set(newValue T) {
	self.previous = self.value
	self.value = Some(newValue)
}

// Value provides the current value.
func (self *WithPrevious[T]) Value() Option[T] {
	return self.value
}
