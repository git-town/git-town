package mapstools

import (
	"cmp"
	"iter"
	"maps"
	"slices"
)

func SortedKeys[Map ~map[K]V, K cmp.Ordered, V any](data Map) iter.Seq[K] {
	return func(yield func(K) bool) {
		for _, key := range slices.Sorted(maps.Keys(data)) {
			if !yield(key) {
				return
			}
		}
	}
}
