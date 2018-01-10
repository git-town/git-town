package cmd

import (
	"fmt"

	"github.com/Originate/git-town/src/prompt"
	"github.com/Originate/git-town/src/script"
)

// These variables represent command-line flags
var (
	abortFlag,
	allFlag,
	continueFlag,
	debugFlag,
	dryRunFlag,
	globalFlag,
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

func validateBooleanArgument(arg string) error {
	if arg != "true" && arg != "false" {
		return fmt.Errorf("Invalid value: '%s'", arg)
	}
	return nil
}

func validateIsConfigured() error {
	prompt.EnsureIsConfigured()
	return nil
}
