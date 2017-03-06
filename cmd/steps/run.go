package steps

import (
  "fmt"
  "os"

  "github.com/Originate/gt/cmd/git"

  "github.com/fatih/color"
)

type RunOptions struct {
  Command string
  IsAbort bool
  IsContinue bool
  IsSkip bool
  SkipMessage string
  StepListGenerator func() StepList
}

func Run(options RunOptions) {
  if options.IsAbort || options.IsContinue {
    runState := loadState(options.Command)
    runState.Command = options.Command
    runState.SkipMessage = options.SkipMessage
    if options.IsAbort {
      abortRunState := runState.CreateAbortRunState()
      runSteps(&abortRunState)
    } else {
      git.EnsureDoesNotHaveConflicts()
      runSteps(&runState)
    }
  } else {
    runSteps(&RunState{
      Command: options.Command,
      RunStepList: options.StepListGenerator(),
    })
  }
}

// Helpers

func runSteps(runState *RunState) {
  for {
    step := runState.RunStepList.Pop()
    if step == nil {
      return
    }
    undoStep := step.CreateUndoStep()
    err := step.Run()
    if err != nil {
      runState.AbortStep = step.CreateAbortStep()
      runState.RunStepList.Prepend(step.CreateContinueStep())
      saveState(runState)
      exitWithMessages(runState.Command, runState.SkipMessage)
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
