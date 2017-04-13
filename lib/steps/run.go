package steps

import (
	"fmt"
	"os"

	"github.com/Originate/git-town/lib/git"
	"github.com/Originate/git-town/lib/util"

	"github.com/fatih/color"
)

// RunOptions bundles the parameters for running a Git Town command.
type RunOptions struct {
	CanSkip              func() bool
	Command              string
	IsAbort              bool
	IsContinue           bool
	IsSkip               bool
	IsUndo               bool
	SkipMessageGenerator func() string
	StepListGenerator    func() StepList
}

// Run runs the Git Town command described by the given RunOptions.
func Run(options RunOptions) {
	if options.IsAbort {
		runState := loadState(options.Command)
		abortRunState := runState.CreateAbortRunState()
		runSteps(&abortRunState, options)
	} else if options.IsContinue {
		runState := loadState(options.Command)
		git.EnsureDoesNotHaveConflicts()
		runSteps(&runState, options)
	} else if options.IsSkip {
		runState := loadState(options.Command)
		skipRunState := runState.CreateSkipRunState()
		runSteps(&skipRunState, options)
	} else if options.IsUndo {
		runState := loadState(options.Command)
		undoRunState := runState.CreateUndoRunState()
		if undoRunState.RunStepList.isEmpty() {
			util.ExitWithErrorMessage("Nothing to undo")
		} else {
			runSteps(&undoRunState, options)
		}
	} else {
		clearSavedState(options.Command)
		runSteps(&RunState{
			Command:     options.Command,
			RunStepList: options.StepListGenerator(),
		}, options)
	}
}

// Helpers

func runSteps(runState *RunState, options RunOptions) {
	for {
		step := runState.RunStepList.Pop()
		if step == nil {
			if !runState.IsAbort && !runState.isUndo {
				runState.AbortStep = NoOpStep{}
				saveState(runState)
			}
			fmt.Println()
			return
		}
		if getTypeName(step) == "SkipCurrentBranchSteps" {
			runState.SkipCurrentBranchSteps()
			continue
		}
		if getTypeName(step) == "PushBranchAfterCurrentBranchSteps" {
			runState.AddPushBranchStepAfterCurrentBranchSteps()
			continue
		}
		undoStepBeforeRun := step.CreateUndoStepBeforeRun()
		err := step.Run()
		if err != nil {
			runState.AbortStep = step.CreateAbortStep()
			if step.ShouldAutomaticallyAbortOnError() {
				abortRunState := runState.CreateAbortRunState()
				runSteps(&abortRunState, options)
				util.ExitWithErrorMessage(step.GetAutomaticAbortErrorMessage())
			} else {
				runState.RunStepList.Prepend(step.CreateContinueStep())
				saveState(runState)
				skipMessage := ""
				if options.CanSkip() {
					skipMessage = options.SkipMessageGenerator()
				}
				exitWithMessages(runState.Command, skipMessage)
			}
		}
		undoStepAfterRun := step.CreateUndoStepAfterRun()
		runState.UndoStepList.Prepend(undoStepBeforeRun)
		runState.UndoStepList.Prepend(undoStepAfterRun)
	}
}

func exitWithMessages(command string, skipMessage string) {
	messageFmt := color.New(color.FgRed)
	fmt.Println()
	messageFmt.Printf("To abort, run \"gt %s --abort\".", command)
	fmt.Println()
	messageFmt.Printf("To continue after you have resolved the conflicts, run \"gt %s --continue\".", command)
	fmt.Println()
	if skipMessage != "" {
		messageFmt.Printf("To skip %s, run \"gt %s --skip\".", skipMessage, command)
		fmt.Println()
	}
	fmt.Println()
	os.Exit(1)
}
