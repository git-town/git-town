package slice

import "slices"

// AppendAllMissing provides the given list with all missing elements of `additional` appended.
func AppendAllMissing[S ~[]C, C comparable](existing S, additional ...C) S { //nolint:ireturn
	result := existing
	for a := range additional {
		if !slices.Contains(result, additional[a]) {
			result = append(result, additional[a])
		}
	}
	return result
}
