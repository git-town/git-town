package asserts

func NoError(err error) {
	if err != nil {
		panic(err)
	}
}

// NoError1 verifies the given result consisting of a value and error.
// If the error is nil, it returns the value, otherwise it panics.
func NoError1[T any](value T, err error) T { //nolint:ireturn
	NoError(err)
	return value
}

// NoError1 verifies the given result consisting of two values and error.
// If the error is nil, it returns the values, otherwise it panics.
func NoError2[A, B any](value1 A, value2 B, err error) (A, B) { //nolint:ireturn
	NoError(err)
	return value1, value2
}
