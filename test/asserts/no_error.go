package asserts

func NoError(err error) {
	if err != nil {
		panic(err)
	}
}

// transparently checks the error value returned by a fallible operation
func NoError1[T any](value T, err error) T { //nolint:ireturn
	NoError(err)
	return value
}

// transparently checks the error value returned by a fallible operation
func NoError2[A, B any](value1 A, value2 B, err error) (A, B) { //nolint:ireturn
	NoError(err)
	return value1, value2
}
