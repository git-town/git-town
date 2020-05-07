package steps

import (
	"fmt"
	"os"

	"github.com/git-town/git-town/src/git"
	"github.com/git-town/git-town/src/util"

	"github.com/fatih/color"
)

// Run runs the Git Town command described by the given state
// nolint: gocyclo, gocognit
func Run(runState *RunState) error {
	for {
		step := runState.RunStepList.Pop()
		if step == nil {
			runState.MarkAsFinished()
			if runState.IsAbort || runState.isUndo {
				err := DeletePreviousRunState()
				if err != nil {
					return fmt.Errorf("cannot delete previous run state: %w", err)
				}
			} else {
				err := SaveRunState(runState)
				if err != nil {
					return fmt.Errorf("cannot save run state: %w", err)
				}
			}
			fmt.Println()
			return nil
		}
		if getTypeName(step) == "*SkipCurrentBranchSteps" {
			runState.SkipCurrentBranchSteps()
			continue
		}
		if getTypeName(step) == "*PushBranchAfterCurrentBranchSteps" {
			runState.AddPushBranchStepAfterCurrentBranchSteps()
			continue
		}
		err := step.Run()
		if err != nil {
			runState.AbortStepList.Append(step.CreateAbortStep())
			if step.ShouldAutomaticallyAbortOnError() {
				abortRunState := runState.CreateAbortRunState()
				err := Run(&abortRunState)
				if err != nil {
					return fmt.Errorf("cannot run the abort steps: %w", err)
				}
				util.ExitWithErrorMessage(step.GetAutomaticAbortErrorMessage())
			} else {
				runState.RunStepList.Prepend(step.CreateContinueStep())
				runState.MarkAsUnfinished()
				if runState.Command == "sync" && !(git.IsRebaseInProgress() && git.Config().IsMainBranch(git.GetCurrentBranchName())) {
					runState.UnfinishedDetails.CanSkip = true
				}
				err := SaveRunState(runState)
				if err != nil {
					return fmt.Errorf("cannot save run state: %w", err)
				}
				exitWithMessages(runState.UnfinishedDetails.CanSkip)
			}
		}
		runState.UndoStepList.Prepend(step.CreateUndoStep())
	}
}

// Helpers

func exitWithMessages(canSkip bool) {
	messageFmt := color.New(color.FgRed)
	fmt.Println()
	util.PrintlnColor(messageFmt, "To abort, run \"git-town abort\".")
	util.PrintlnColor(messageFmt, "To continue after having resolved conflicts, run \"git-town continue\".")
	if canSkip {
		util.PrintlnColor(messageFmt, "To continue by skipping the current branch, run \"git-town skip\".")
	}
	fmt.Println()
	os.Exit(1)
}
