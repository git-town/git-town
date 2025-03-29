package mapstools

import (
	"cmp"
	"maps"
	"slices"
)

// SortedKeyValues provides the keys and values of the given map, sorted alphabetically by key.
func SortedKeyValues[K cmp.Ordered, V any](data map[K]V) []struct {
	Key   K
	Value V
} {
	result := make([]struct {
		Key   K
		Value V
	}, len(data))
	for k, key := range slices.Sorted(maps.Keys(data)) {
		result[k] = struct {
			Key   K
			Value V
		}{
			Key:   key,
			Value: data[key],
		}
	}
	return result
}
