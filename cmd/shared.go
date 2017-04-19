package cmd

import "errors"

// These variables represent command-line flags
var (
	abortFlag,
	allFlag,
	continueFlag,
	skipFlag,
	undoFlag bool
)

func validateArgsCount(args []string, count int) error {
	err := validateMinArgs(args, count)
	if err != nil {
		return err
	}
	return validateMaxArgs(args, count)
}

func validateMinArgs(args []string, min int) error {
	if len(args) < min {
		return errors.New("Too few arguments")
	}
	return nil
}

func validateMaxArgs(args []string, max int) error {
	if len(args) > max {
		return errors.New("Too many arguments")
	}
	return nil
}
