package helpers

import (
	"fmt"
	"strings"
)

// OrderedSet is a Set that provides its elements in the order they were received.
type OrderedSet[T comparable] struct {
	elements []T
}

// NewOrderedSet provides instances of OrderedSet populated with the given elements.
func NewOrderedSet[T comparable](elements ...T) OrderedSet[T] {
	return OrderedSet[T]{elements}
}

// Add provides a new OrderedSet with the given element added.
// The element is only added if it doesn't exist in the original set.
// TODO: rename os to self
func (os OrderedSet[T]) Add(element T) OrderedSet[T] {
	if !os.Contains(element) {
		return OrderedSet[T]{append(os.elements, element)}
	}
	return os
}

// Contains indicates whether this Set contains the given element.
func (os OrderedSet[T]) Contains(element T) bool {
	for _, existing := range os.elements {
		if element == existing {
			return true
		}
	}
	return false
}

// Elements provides the elements of this os in the order they were received.
func (os OrderedSet[T]) Elements() []T {
	return os.elements
}

func (os OrderedSet[T]) Join(sep string) string {
	texts := make([]string, len(os.elements))
	for e, element := range os.elements {
		texts[e] = fmt.Sprintf("%v", element)
	}
	return strings.Join(texts, sep)
}
