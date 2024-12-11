package asserts

// transparently checks the error value returned by a fallible operation
func NoError1[T any](value T, err error) T { //nolint:ireturn
	NoError(err)
	return value
}
