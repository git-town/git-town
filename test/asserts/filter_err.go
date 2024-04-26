package asserts

func FilterErr[T any](data T, err error) T { //nolint:ireturn
	if err != nil {
		panic(err.Error())
	}
	return data
}
