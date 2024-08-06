package mapstools

import "maps"

// provides the given maps merged together
func Merge[K comparable, V any](map1, map2 map[K]V) map[K]V {
	result := maps.Clone(map1)
	maps.Copy(result, map2)
	return result
}
