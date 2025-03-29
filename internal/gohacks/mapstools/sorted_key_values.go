package mapstools

import (
	"cmp"
	"iter"
	"maps"
	"slices"
)

func SortedKeyValues[Map ~map[K]V, K cmp.Ordered, V any](data Map) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for _, key := range slices.Sorted(maps.Keys(data)) {
			if !yield(key, data[key]) {
				return
			}
		}
	}
}
