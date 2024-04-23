package option

// Unpack
func Unpack[T any](option *T, err error) (T, error) {
	if option == nil {
		var emptyT T
		return emptyT, err
	}
	return *option, nil
}
