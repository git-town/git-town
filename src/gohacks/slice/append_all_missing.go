package slice

// AppendAllMissing provides the given list with all missing elements of `additional` appended.
func AppendAllMissing[S ~[]C, C comparable](existing S, additional ...C) S { //nolint:ireturn
	result := existing
	for a := range additional {
		if !Contains(result, additional[a]) {
			result = append(result, additional[a])
		}
	}
	return result
}
