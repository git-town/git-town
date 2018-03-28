package steps

import (
	"fmt"
	"os"

	"github.com/Originate/exit"
	"github.com/Originate/git-town/src/git"
	"github.com/Originate/git-town/src/util"

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
		runState := LoadPreviousRunState(options.Command)
		abortRunState := runState.CreateAbortRunState()
		runSteps(&abortRunState, options)
	} else if options.IsContinue {
		runState := LoadPreviousRunState(options.Command)
		if runState.RunStepList.isEmpty() {
			util.ExitWithErrorMessage("Nothing to continue")
		}
		git.EnsureDoesNotHaveConflicts()
		runSteps(runState, options)
	} else if options.IsSkip {
		runState := LoadPreviousRunState(options.Command)
		skipRunState := runState.CreateSkipRunState()
		runSteps(&skipRunState, options)
	} else if options.IsUndo {
		runState := LoadPreviousRunState(options.Command)
		undoRunState := runState.CreateUndoRunState()
		if undoRunState.RunStepList.isEmpty() {
			util.ExitWithErrorMessage("Nothing to undo")
		} else {
			runSteps(&undoRunState, options)
		}
	} else {
		DeleteRunState(options.Command)
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
				runSteps(&abortRunState, options)
				util.ExitWithErrorMessage(step.GetAutomaticAbortErrorMessage())
			} else {
				runState.RunStepList.Prepend(step.CreateContinueStep())
				runState.Save()
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
	_, err := messageFmt.Printf("To abort, run \"git-town %s --abort\".\n", command)
	exit.If(err)
	_, err = messageFmt.Printf("To continue after you have resolved the conflicts, run \"git-town %s --continue\".\n", command)
	exit.If(err)
	if skipMessage != "" {
		_, err = messageFmt.Printf("To skip %s, run \"git-town %s --skip\".\n", skipMessage, command)
		exit.If(err)
	}
	fmt.Println()
	os.Exit(1)
}
