package steps

import (
  "fmt"
  "strings"

  "github.com/Originate/gt/lib/git"
)


type Step interface {
  CreateAbortStep() Step
  CreateContinueStep() Step
  CreateUndoStep() Step
  Run() error
}


type SerializedStep struct {
  Data []byte
  Type string
}


type SerializedRunState struct {
  AbortStep SerializedStep
  RunSteps []SerializedStep
  UndoSteps []SerializedStep
}

func getRunResultFilename(command string) string {
  return fmt.Sprintf("/tmp/%s_%s", command, strings.Replace(git.GetRootDirectory(), "/", "_", -1))
}
