package util

// FirstError collects errors
func FirstError(errors ...interface{}) error {
	for _, err := range errors {
		if err == nil {
			continue
		}
		switch err.(type) {
		case error:
			return err.(error)
		case func() error:
			f := err.(func() error)
			result := f()
			if result != nil {
				return result
			}
		case func():
			f := err.(func())
			f()
		default:
			panic("Unknown type provided to util.FirstError")
		}
	}
	return nil
}
