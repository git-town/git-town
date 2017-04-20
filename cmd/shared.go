package cmd

import (
	"errors"
	"fmt"
	"log"
	"strconv"
)

// These variables represent command-line flags
var (
	abortFlag,
	allFlag,
	continueFlag,
	skipFlag,
	undoFlag bool
)

func stringToBool(arg string) bool {
	value, err := strconv.ParseBool(arg)
	if err != nil {
		log.Fatal(err)
	}
	return value
}

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
