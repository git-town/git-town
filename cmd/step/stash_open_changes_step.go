package step

import (
  "github.com/Originate/gt/cmd/script"
)

type StashOpenChangesStep int

func (step StashOpenChangesStep) CreateAbortStep() Step {
  return new(NoOpStep)
}

func (step StashOpenChangesStep) CreateContinueStep() Step {
  return new(NoOpStep)
}

func (step StashOpenChangesStep) CreateUndoStep() Step {
  return new(NoOpStep)
}

func (step StashOpenChangesStep) Run() error {
  err := script.RunCommand([]string{"git", "add", "-A"})
  if err != nil {
    return err
  }
  return script.RunCommand([]string{"git", "stash"})
}
