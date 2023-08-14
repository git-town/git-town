package genericslice

func AppendAllMissing[C comparable](list []C, elements []C) []C {
	for _, element := range elements {
		if !Contains(list, element) {
			list = append(list, element)
		}
	}
	return list
}

// Contains returns whether the given string slice
// contains the given string.
func Contains[C comparable](list []C, value C) bool {
	for _, element := range list {
		if element == value {
			return true
		}
	}
	return false
}

// FirstElementOr provides the first element of the given list or the given alternative if the list is empty.
func FirstElementOr[C comparable](list []C, alternative C) C {
	if len(list) > 0 {
		return list[0]
	}
	return alternative
}

// Hoist provides the given list of strings, with the given element moved to the first position.
func Hoist[C comparable](list []C, element C) []C {
	result := make([]C, 0, len(list))
	hasElement := false
	for _, input := range list {
		if input == element {
			hasElement = true
		} else {
			result = append(result, input)
		}
	}
	if hasElement {
		result = append([]C{element}, result...)
	}
	return result
}

// Remove returns a new string slice which is the given string slice
// with the given string removed.
func Remove[C comparable](list []C, value C) []C {
	result := make([]C, 0, len(list)-1)
	for _, element := range list {
		if element != value {
			result = append(result, element)
		}
	}
	return result
}
