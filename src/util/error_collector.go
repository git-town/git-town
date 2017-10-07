package util

// FirstError collects errors
func FirstError(errorFuncs ...func() error) error {
	for _, errorFunc := range errorFuncs {
		err := errorFunc()
		if err != nil {
			return err
		}
	}
	return nil
}
