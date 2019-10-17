package helpers

import "strings"

// OrderedStringSet is a Set for strings that provides the accumulated strings in the order they were received.
// The zero value is a fully functional empty set.
type OrderedStringSet struct {
	elements []string
}

// Add appends the given text to this set.
// If the text already exists, it re-uses the existing element and does not append a new one.
func (set *OrderedStringSet) Add(text string) {
	if !set.Contains(text) {
		set.elements = append(set.elements, text)
	}
}

// Contains indicates whether this Set contains the given string.
func (set *OrderedStringSet) Contains(text string) bool {
	for _, element := range set.elements {
		if text == element {
			return true
		}
	}
	return false
}

// Slice provides the elements of this set in the order they were received.
func (set *OrderedStringSet) Slice() []string {
	return set.elements
}

func (set *OrderedStringSet) String() string {
	return strings.Join(set.elements, ", ")
}
