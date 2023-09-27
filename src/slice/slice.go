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

// LastIndex provides the zero-based index of the last occurrence of the given element in the given list.
func LastIndex[C comparable](list []C, element C) int {
	for l := len(list) - 1; l >= 0; l-- {
		if list[l] == element {
			return l
		}
	}
	return -1
}

// LowerLast provides the given slice with the last element of the given type moved to the last position in the list.
func LowerLast[C comparable](list []C, element C) []C {
	index := LastIndex(list, element)
	if index == -1 {
		return list
	}
	removed := RemoveAt(list, index)
	return append(removed, element)
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

// RemoveAt provides the given list with the element at the given position removed.
func RemoveAt[C comparable](list []C, index int) []C {
	return append(list[:index], list[index+1:]...)
}
