package step

import (
  "fmt"
  "strings"

  "github.com/Originate/gt/cmd/git"
)

type Step interface {
  CreateAbortStep() Step
  CreateContinueStep() Step
  CreateUndoStep() Step
  Run() error
}

type RunResult struct {
  AbortSteps []Step
  ContinueSteps []Step
  Success bool
  UndoSteps []Step
}

type SerializedStep struct {
  Data []byte
  Type string
}

type SerializedRunResult struct {
  AbortSteps []SerializedStep
  ContinueSteps []SerializedStep
  Success bool
  UndoSteps []SerializedStep
}

func getRunResultFilename(commandName string) string {
  return fmt.Sprintf("/tmp/%s_%s", commandName, strings.Replace(git.GetRootDirectory(), "/", "_", -1))
}
