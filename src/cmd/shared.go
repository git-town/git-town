package cmd

import (
	"fmt"
	"os"

	"github.com/Originate/git-town/src/git"
	"github.com/Originate/git-town/src/prompt"
	"github.com/Originate/git-town/src/script"
	"github.com/Originate/git-town/src/steps"
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

func ensureIsNotInUnfinishedState() error {
	runState := steps.LoadPreviousRunState()
	if runState != nil && runState.IsUnfinished() {
		response := prompt.AskHowToHandleUnfinishedRunState(
			runState.Command,
			runState.UnfinishedDetails.EndBranch,
			runState.UnfinishedDetails.EndTime,
			runState.UnfinishedDetails.CanSkip,
		)
		if response == prompt.ResponseTypeDiscard {
			steps.DeletePreviousRunState()
			return nil
		}
		switch response {
		case prompt.ResponseTypeContinue:
			git.EnsureDoesNotHaveConflicts()
			steps.Run(runState)
		case prompt.ResponseTypeAbort:
			abortRunState := runState.CreateAbortRunState()
			steps.Run(&abortRunState)
		case prompt.ResponseTypeSkip:
			skipRunState := runState.CreateSkipRunState()
			steps.Run(&skipRunState)
		}
		os.Exit(0)
	}
	return nil
}
