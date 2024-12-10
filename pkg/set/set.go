package set

import (
	"cmp"
	"slices"

	"golang.org/x/exp/maps"
)

// Set implements a simple generic Set implementation.
// The zero value is not valid. Use set.New to create new instances.
type Set[T cmp.Ordered] map[T]struct{}

func New[T cmp.Ordered](values ...T) Set[T] {
	result := Set[T]{}
	result.Add(values...)
	return result
}

func (self Set[T]) Add(values ...T) {
	for _, value := range values {
		self[value] = struct{}{}
	}
}

func (self Set[T]) AddSet(other Set[T]) {
	self.Add(other.Values()...)
}

func (self Set[T]) Contains(value T) bool {
	_, has := self[value]
	return has
}

func (self Set[T]) Values() []T {
	result := maps.Keys(self)
	slices.Sort(result)
	return result
}
