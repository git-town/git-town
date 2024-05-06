package slice

import "slices"

// RemoveAt provides the given list without the elements at the given positions.
func RemoveAt[S ~[]C, C comparable](list S, indexes ...int) S { //nolint: ireturn
	result := make(S, 0, len(list)-len(indexes))
	for l := range list {
		if !slices.Contains(indexes, l) {
			result = append(result, list[l])
		}
	}
	return result
}
