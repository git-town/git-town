package step

import (
  "github.com/Originate/gt/cmd/script"
)

type StashOpenChangesStep struct {}

func (step StashOpenChangesStep) CreateAbortStep() Step {
  return nil
}

func (step StashOpenChangesStep) CreateContinueStep() Step {
  return nil
}

func (step StashOpenChangesStep) CreateUndoStep() Step {
  return RestoreOpenChangesStep{}
}

func (step StashOpenChangesStep) Run() error {
  err := script.RunCommand([]string{"git", "add", "-A"})
  if err != nil {
    return err
  }
  return script.RunCommand([]string{"git", "stash"})
}
