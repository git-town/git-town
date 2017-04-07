package cmd

import "errors"

func validateMaxArgs(args []string, max int) error {
	if len(args) > max {
		return errors.New("Too many arguments")
	}
	return nil
}
