package mapstools

import (
	"cmp"
	"iter"
)

// SortedKeyValues provides the keys and values of the given map sorted alphabetically by key.
func SortedKeyValues[Map ~map[K]V, K cmp.Ordered, V any](data Map) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for key := range SortedKeys(data) {
			if !yield(key, data[key]) {
				return
			}
		}
	}
}
