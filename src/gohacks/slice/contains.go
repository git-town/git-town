package slice

// Contains indicates whether the given slice contains the given element.
func Contains[C comparable](list []C, value C) bool {
	for l := range list {
		if list[l] == value {
			return true
		}
	}
	return false
}
