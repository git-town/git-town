package steps

import (
	"fmt"
	"os"

	"github.com/Originate/exit"
	"github.com/Originate/git-town/src/git"
	"github.com/Originate/git-town/src/util"

	"github.com/fatih/color"
)

// Run runs the Git Town command described by the given RunOptions.
// func Run(options RunOptions) {
// 	if options.IsAbort {
// 		runState := LoadPreviousRunState(options.Command)
// 		abortRunState := runState.CreateAbortRunState()
// 		runSteps(&abortRunState, options)
// 	} else if options.IsContinue {
// 		runState := LoadPreviousRunState(options.Command)
// 		if runState.RunStepList.isEmpty() {
// 			util.ExitWithErrorMessage("Nothing to continue")
// 		}
// 		git.EnsureDoesNotHaveConflicts()
// 		runSteps(runState, options)
// 	} else if options.IsSkip {
// 		runState := LoadPreviousRunState(options.Command)
// 		skipRunState := runState.CreateSkipRunState()
// 		runSteps(&skipRunState, options)
// 	} else if options.IsUndo {
// 		runState := LoadPreviousRunState(options.Command)
// 		undoRunState := runState.CreateUndoRunState()
// 		if undoRunState.RunStepList.isEmpty() {
// 			util.ExitWithErrorMessage("Nothing to undo")
// 		} else {
// 			runSteps(&undoRunState, options)
// 		}
// 	} else {
// 		DeletePreviousRunState(options.Command)
// 		runSteps(&RunState{
// 			Command:     options.Command,
// 			RunStepList: options.StepListGenerator(),
// 		}, options)
// 	}
// }

// Helpers

func Run(runState *RunState) {
	for {
		step := runState.RunStepList.Pop()
		if step == nil {
			runState.MarkAsFinished()
			if !runState.IsAbort && !runState.isUndo {
				runState.Save()
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
				skipMessage := ""
				if runState.Command == "sync" && !(git.IsRebaseInProgress() && git.IsMainBranch(git.GetCurrentBranchName())) {
					runState.CanSkip = true
					skipMessage = fmt.Sprintf("the sync of the '%s' branch", git.GetCurrentBranchName())
				}
				runState.Save()
				exitWithMessages(skipMessage)
			}
		}
		undoStepAfterRun := step.CreateUndoStepAfterRun()
		runState.UndoStepList.Prepend(undoStepBeforeRun)
		runState.UndoStepList.Prepend(undoStepAfterRun)
	}
}

func exitWithMessages(skipMessage string) {
	messageFmt := color.New(color.FgRed)
	fmt.Println()
	_, err := messageFmt.Printf("To abort, run \"git-town abort\".\n")
	exit.If(err)
	_, err = messageFmt.Printf("To continue after you have resolved the conflicts, run \"git-town continue\".\n")
	exit.If(err)
	if skipMessage != "" {
		_, err = messageFmt.Printf("To skip %s, run \"git-town skip\".\n", skipMessage)
		exit.If(err)
	}
	fmt.Println()
	os.Exit(1)
}
