package slice

import (
	"slices"
)

// FindAllMissing provides all entries in additional that are not contained in existing.
func FindAllMissing[S ~[]C, C comparable](existing, additional S) S { //nolint:ireturn
	result := S{}
	for _, element := range additional {
		if !slices.Contains(existing, element) {
			result = append(result, element)
		}
	}
	return result
}
