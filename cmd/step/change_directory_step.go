package step

import (
  "github.com/Originate/gt/cmd/script"
)

type ChangeDirectoryStep struct {
  directory string
}

func (step ChangeDirectoryStep) CreateAbortStep() Step {
  return new(NoOpStep)
}

func (step ChangeDirectoryStep) CreateContinueStep() Step {
  return new(NoOpStep)
}

func (step ChangeDirectoryStep) CreateUndoStep() Step {
  return new(NoOpStep)
}

func (step ChangeDirectoryStep) Run() error {
  return script.RunCommand([]string{"cd", step.directory})
}
