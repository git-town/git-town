package gohacks

import "golang.org/x/exp/maps"

// a simple generic Set implementation
type Set[T comparable] map[T]struct{}

func NewSet[T comparable](values ...T) Set[T] {
	result := Set[T]{}
	for _, value := range values {
		result.Add(value)
	}
	return result
}

func (self Set[T]) Add(value T) {
	self[value] = struct{}{}
}

func (self Set[T]) AddMany(values ...T) {
	for _, value := range values {
		self.Add(value)
	}
}

func (self Set[T]) AddSet(other Set[T]) {
	for _, value := range other.Values() {
		self.Add(value)
	}
}

func (self Set[T]) Contains(value T) bool {
	_, has := self[value]
	return has
}

func (self Set[T]) Values() []T {
	return maps.Keys(self)
}
