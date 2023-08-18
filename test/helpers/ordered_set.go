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
func (set OrderedSet[T]) Add(element T) OrderedSet[T] {
	if !set.Contains(element) {
		return OrderedSet[T]{append(set.elements, element)}
	}
	return set
}

// Contains indicates whether this Set contains the given element.
func (set OrderedSet[T]) Contains(element T) bool {
	for _, existing := range set.elements {
		if element == existing {
			return true
		}
	}
	return false
}

// Elements provides the elements of this set in the order they were received.
func (set OrderedSet[T]) Elements() []T {
	return set.elements
}

func (set OrderedSet[T]) Join(sep string) string {
	texts := []string{}
	for _, element := range set.elements {
		texts = append(texts, fmt.Sprintf("%v", element)) // TODO: this might not work as intended
	}
	return strings.Join(texts, sep)
}
