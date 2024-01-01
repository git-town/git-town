package slice

// FirstElementOr provides the first element of the given list or the given alternative if the list is empty.
func FirstElementOr[C comparable](list []C, alternative C) C {
	if len(list) > 0 {
		return list[0]
	}
	return alternative
}
