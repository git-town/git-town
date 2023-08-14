package slice

// AppendAllMissing appends all elements of `additional` that aren't contained in `existing` to `existing`.
func AppendAllMissing[C comparable](existing []C, additional []C) []C {
	for a := range additional {
		if !Contains(existing, additional[a]) {
			existing = append(existing, additional[a])
		}
	}
	return existing
}

// Contains returns whether the given slice contains the given element.
func Contains[C comparable](list []C, value C) bool {
	for l := range list {
		if list[l] == value {
			return true
		}
	}
	return false
}

// FirstElementOr provides the first element of the given list or the given alternative if the list is empty.
func FirstElementOr[C comparable](list []C, alternative C) C { //nolint:ireturn // there should never be any nil values here
	if len(list) > 0 {
		return list[0]
	}
	return alternative
}

// Hoist provides the given slice with the given element moved to the first position.
func Hoist[C comparable](list []C, element C) []C {
	result := make([]C, 0, len(list))
	hasElement := false
	for l := range list {
		if list[l] == element {
			hasElement = true
		} else {
			result = append(result, list[l])
		}
	}
	if !hasElement {
		return result
	}
	return append([]C{element}, result...)
}

// Remove returns a new slice which is the given slice with the given element removed.
func Remove[C comparable](list []C, value C) []C {
	result := make([]C, 0, len(list)-1)
	for l := range list {
		if list[l] != value {
			result = append(result, list[l])
		}
	}
	return result
}
