package slice

import (
	"slices"
)

// FindAllMissing provides all entries in additional that are not contained in existing.
func FindAllMissing[S ~[]C, C comparable](existing, additional S) S { //nolint:ireturn
	result := S{}
	for a := range additional {
		if !slices.Contains(existing, additional[a]) {
			result = append(result, additional[a])
		}
	}
	return result
}
