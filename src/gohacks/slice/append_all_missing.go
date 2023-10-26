package slice

// AppendAllMissing appends all elements of `additional` that aren't contained in `existing` to `existing`.
func AppendAllMissing[S ~[]C, C comparable](existing S, additional S) S { //nolint:ireturn
	for a := range additional {
		if !Contains(existing, additional[a]) {
			existing = append(existing, additional[a])
		}
	}
	return existing
}
