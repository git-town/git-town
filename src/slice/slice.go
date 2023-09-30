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

// Contains returns whether the given slice contains the given element.
func Contains[C comparable](list []C, value C) bool {
	for l := range list {
		if list[l] == value {
			return true
		}
	}
	return false
}

// FindAll provides the indexes of all occurrences of the given element in the given list.
func FindAll[C comparable](list []C, element C) []int {
	result := []int{}
	for l, li := range list {
		if li == element {
			result = append(result, l)
		}
	}
	return result
}

// FirstElementOr provides the first element of the given list or the given alternative if the list is empty.
func FirstElementOr[C comparable](list []C, alternative C) C { //nolint:ireturn
	if len(list) > 0 {
		return list[0]
	}
	return alternative
}

// Hoist provides the given slice with the given element moved to the first position.
func Hoist[S ~[]C, C comparable](list S, element C) S { //nolint:ireturn
	result := make([]C, 0, len(list))
	hasElement := false
	for l := range list {
		if list[l] == element {
			hasElement = true
		} else {
			result = append(result, list[l])
		}
	}
	// TODO: remove the negation from the if condition so that there is only one return statement
	if !hasElement {
		return result
	}
	return append([]C{element}, result...)
}

// LowerLast provides the given slice with the last element of the given type moved to the last position in the list.
func LowerAll[C comparable](haystack []C, needle C) []C {
	result := make([]C, 0, len(haystack))
	hasNeedle := false
	for _, element := range haystack {
		if element == needle {
			hasNeedle = true
		} else {
			result = append(result, element)
		}
	}
	if hasNeedle {
		result = append(result, needle)
	}
	return result
}

// Remove returns a new slice which is the given slice with the given element removed.
func Remove[S ~[]C, C comparable](list S, value C) S { //nolint:ireturn
	result := make([]C, 0, len(list)-1)
	for l := range list {
		if list[l] != value {
			result = append(result, list[l])
		}
	}
	return result
}

// RemoveAt provides the given list with the element at the given position removed.
func RemoveAt[S ~[]C, C comparable](list S, index int) S { //nolint:ireturn
	return append(list[:index], list[index+1:]...)
}

// TruncateLast provides the given list with its last element removed.
func TruncateLast[S ~[]C, C comparable](list S) S { //nolint:ireturn
	listLength := len(list)
	if listLength == 0 {
		return list
	}
	return list[:listLength-1]
}
