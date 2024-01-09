package slice

// AppendAllMissing appends all elements of `additional` that aren't contained in `existing` to `existing`.
func AppendAllMissing[S ~[]C, C comparable](existing S, additional S) S { //nolint:ireturn
	result := existing
	for a := range additional {
		if !Contains(result, additional[a]) {
			result = append(result, additional[a])
		}
	}
	return result
}
