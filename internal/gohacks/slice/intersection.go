package slice

import (
	"slices"
)

func Intersection[S ~[]C, C comparable](existing, additional S) S { //nolint:ireturn
	result := S{}
	for a := range additional {
		if !slices.Contains(existing, additional[a]) {
			result = append(result, additional[a])
		}
	}
	return result
}
