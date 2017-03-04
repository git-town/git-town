package stepRunner

import (
  "fmt"
  "os"

  "github.com/Originate/gt/cmd/git"
  "github.com/Originate/gt/cmd/step"

  "github.com/fatih/color"
)

type Options struct {
  Command string
  IsAbort bool
  IsContinue bool
  IsSkip bool
  SkipMessage string
  StepGenerator func() []step.Step
}

func Run(options Options) {
  if options.IsAbort || options.IsContinue {
    runResult := step.Import(options.Command)
    if options.IsAbort {
      steps := append([]step.Step{runResult.AbortStep}, runResult.UndoSteps...)
      runSteps(steps, []step.Step{}, options.Command, options.SkipMessage)
    } else {
      git.EnsureDoesNotHaveConflicts()
      runSteps(runResult.ContinueSteps, runResult.UndoSteps, options.Command, options.SkipMessage)
    }
  } else {
    steps := options.StepGenerator()
    runSteps(steps, []step.Step{}, options.Command, options.SkipMessage)
  }
}

// Helpers

func runSteps(steps, undoSteps []step.Step, command string, skipMessage string) {
  for i := 0; i < len(steps); i++ {
    undoStep := steps[i].CreateUndoStep()
    err := steps[i].Run()
    if err != nil {
      abortStep := steps[i].CreateAbortStep()
      continueSteps := append([]step.Step{steps[i].CreateContinueStep()}, steps[i+1:]...)
      step.Export(command, abortStep, continueSteps, undoSteps)
      exitWithMessages(command, skipMessage)
    }
    undoSteps = append([]step.Step{undoStep}, undoSteps...)
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
