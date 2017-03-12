package steps

import (
  "fmt"
  "os"

  "github.com/Originate/gt/lib/git"

  "github.com/fatih/color"
)

type RunOptions struct {
  CanSkip func() bool
  Command string
  IsAbort bool
  IsContinue bool
  IsSkip bool
  SkipMessageGenerator func() string
  StepListGenerator func() StepList
}

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
  } else {
    runSteps(&RunState{
      Command: options.Command,
      RunStepList: options.StepListGenerator(),
    }, options)
  }
}

// Helpers

func runSteps(runState *RunState, options RunOptions) {
  for {
    step := runState.RunStepList.Pop()
    if step == nil {
      return
    }
    if getTypeName(step) == "SkipCurrentBranchSteps" {
      runState.SkipCurrentBranchSteps()
      continue
    }
    undoStep := step.CreateUndoStep()
    err := step.Run()
    if err != nil {
      runState.AbortStep = step.CreateAbortStep()
      runState.RunStepList.Prepend(step.CreateContinueStep())
      saveState(runState)
      skipMessage := ""
      if options.CanSkip() {
        skipMessage = options.SkipMessageGenerator()
      }
      exitWithMessages(runState.Command, skipMessage)
    }
    runState.UndoStepList.Prepend(undoStep)
  }
  fmt.Println()
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
