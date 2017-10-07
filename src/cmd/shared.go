package cmd

import (
	"errors"
	"fmt"

	"github.com/Originate/git-town/src/script"
	"github.com/Originate/git-town/src/util"
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

func conditionallyActivateDryRun() error {
	if dryRunFlag {
		script.ActivateDryRun()
	}
	return nil
}

func validateArgsCountFunc(args []string, count int) func() error {
	return func() error {
		return util.FirstError(
			validateMinArgsFunc(args, count),
			validateMaxArgsFunc(args, count),
		)
	}
}

func validateBooleanArgument(arg string) error {
	if arg != "true" && arg != "false" {
		return fmt.Errorf("Invalid value: '%s'", arg)
	}
	return nil
}

func validateBooleanArgumentFunc(arg string) func() error {
	return func() error {
		return validateBooleanArgument(arg)
	}
}

func validateMinArgs(args []string, min int) error {
	if len(args) < min {
		return errors.New("Too few arguments")
	}
	return nil
}

func validateMinArgsFunc(args []string, min int) func() error {
	return func() error {
		return validateMinArgs(args, min)
	}
}

func validateMaxArgs(args []string, max int) error {
	if len(args) > max {
		return errors.New("Too many arguments")
	}
	return nil
}

func validateMaxArgsFunc(args []string, max int) func() error {
	return func() error {
		return validateMaxArgs(args, max)
	}
}

func validateIsConfigured() error {
	validateIsConfigured()
	return nil
}
