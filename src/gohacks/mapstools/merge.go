package mapstools

import "maps"

func Merge[K comparable, V any](map1, map2 map[K]V) map[K]V {
	result := maps.Clone(map1)
	maps.Copy(result, map2)
	return result
}
