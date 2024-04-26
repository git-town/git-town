package gohacks

func FilterErr[T any](data T, err error) T {
	if err != nil {
		panic(err.Error())
	}
	return data
}
