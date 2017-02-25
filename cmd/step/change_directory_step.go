package step

import (
  "github.com/Originate/gt/cmd/script"
)

type ChangeDirectoryStep struct {
  directory string
}

func (step ChangeDirectoryStep) CreateAbortStep() Step {
  return nil
}

func (step ChangeDirectoryStep) CreateContinueStep() Step {
  return nil
}

func (step ChangeDirectoryStep) CreateUndoStep() Step {
  return nil
}

func (step ChangeDirectoryStep) Run() error {
  return script.RunCommand([]string{"cd", step.directory})
}
