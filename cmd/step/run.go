package step

import (
  "fmt"
  "os"

  "github.com/fatih/color"
)

func Run(steps []Step, undoSteps []Step, command string, skipMessage string) {
  for i := 0; i < len(steps); i++ {
    undoStep := steps[i].CreateUndoStep()
    err := steps[i].Run()
    if err != nil {
      abortStep := steps[i].CreateAbortStep()
      continueSteps := append([]Step{steps[i].CreateContinueStep()}, steps[i+1:]...)
      export(command, abortStep, continueSteps, undoSteps)
      exitWithMessages(command, skipMessage)
    }
    undoSteps = append([]Step{undoStep}, undoSteps...)
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
