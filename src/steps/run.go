package steps

import (
	"fmt"
	"os"

	"github.com/Originate/git-town/src/git"
	"github.com/Originate/git-town/src/util"
	"github.com/pkg/errors"

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
					return errors.Wrap(err, "cannot delete previous run state")
				}
			} else {
				err := SaveRunState(runState)
				if err != nil {
					return errors.Wrap(err, "cannot save run state")
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
		undoStepBeforeRun := step.CreateUndoStepBeforeRun()
		err := step.Run()
		if err != nil {
			runState.AbortStepList.Append(step.CreateAbortStep())
			if step.ShouldAutomaticallyAbortOnError() {
				abortRunState := runState.CreateAbortRunState()
				err := Run(&abortRunState)
				if err != nil {
					return errors.Wrap(err, "cannot run the abort steps")
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
					return errors.Wrap(err, "cannot save run state")
				}
				exitWithMessages(runState.UnfinishedDetails.CanSkip)
			}
		}
		undoStepAfterRun := step.CreateUndoStepAfterRun()
		runState.UndoStepList.Prepend(undoStepBeforeRun)
		runState.UndoStepList.Prepend(undoStepAfterRun)
	}
}

// Helpers

func exitWithMessages(canSkip bool) {
	messageFmt := color.New(color.FgRed)
	fmt.Println()
	// NOTE: No point in handling the errors while printing here,
	//       we are exiting with an error anyways.
	_, _ = messageFmt.Printf("To abort, run \"git-town abort\".\n")
	_, _ = messageFmt.Printf("To continue after having resolved conflicts, run \"git-town continue\".\n")
	if canSkip {
		_, _ = messageFmt.Printf("To continue by skipping the current branch, run \"git-town skip\".\n")
	}
	fmt.Println()
	os.Exit(1)
}
