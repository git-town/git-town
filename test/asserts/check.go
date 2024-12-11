package asserts

func Check[T any](value T, err error) T {
	NoError(err)
	return value
}
