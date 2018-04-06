package steps

import (
	"fmt"
	"os"

	"github.com/Originate/exit"
	"github.com/Originate/git-town/src/git"
	"github.com/Originate/git-town/src/util"

	"github.com/fatih/color"
)

// Run runs the Git Town command described by the given state
// nolint: gocyclo
func Run(runState *RunState) {
	for {
		step := runState.RunStepList.Pop()
		if step == nil {
			runState.MarkAsFinished()
			if !runState.IsAbort && !runState.isUndo {
				SaveRunState(runState)
			}
			fmt.Println()
			return
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
				Run(&abortRunState)
				util.ExitWithErrorMessage(step.GetAutomaticAbortErrorMessage())
			} else {
				runState.RunStepList.Prepend(step.CreateContinueStep())
				runState.MarkAsUnfinished()
				if runState.Command == "sync" && !(git.IsRebaseInProgress() && git.IsMainBranch(git.GetCurrentBranchName())) {
					runState.UnfinishedDetails.CanSkip = true
				}
				SaveRunState(runState)
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
	_, err := messageFmt.Printf("To abort, run \"git-town abort\".\n")
	exit.If(err)
	_, err = messageFmt.Printf("To continue after having resolved conflicts, run \"git-town continue\".\n")
	exit.If(err)
	if canSkip {
		_, err = messageFmt.Printf("To continue by skipping the current branch, run \"git-town skip\".\n")
		exit.If(err)
	}
	fmt.Println()
	os.Exit(1)
}
