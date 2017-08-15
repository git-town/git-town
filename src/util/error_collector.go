package util

// FirstError collects errors
func FirstError(errors ...error) error {
	for _, err := range errors {
		if err != nil {
			return err
		}
	}
	return nil
}

// CollectErrorF collects errors
func CollectErrorF(functions ...func() error) error {
	for _, function := range functions {
		err := function()
		if err != nil {
			return err
		}
	}
	return nil
}
