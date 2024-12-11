package asserts

// transparently checks the error value returned by a fallible operation
func Check1[T any](value T, err error) T { //nolint:ireturn
	NoError(err)
	return value
}
