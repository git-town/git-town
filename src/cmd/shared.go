package cmd

import (
	"errors"
	"fmt"
)

// These variables represent command-line flags
var (
	abortFlag,
	allFlag,
	continueFlag,
	dryRunFlag,
	skipFlag,
	undoFlag bool
)

var abortFlagDescription = "Abort a previous command that resulted in a conflict"
var continueFlagDescription = "Continue a previous command that resulted in a conflict"
var undoFlagDescription = "Undo a previous command"
var dryRunFlagDescription = "Output the commands that would be run without them"

func validateArgsCount(args []string, count int) error {
	err := validateMinArgs(args, count)
	if err != nil {
		return err
	}
	return validateMaxArgs(args, count)
}

func validateBooleanArgument(arg string) error {
	if arg != "true" && arg != "false" {
		return fmt.Errorf("Invalid value: '%s'", arg)
	}
	return nil
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
