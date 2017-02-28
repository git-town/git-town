package step

import (
  "github.com/Originate/gt/cmd/script"
)

type RestoreOpenChangesStep struct {}

func (step RestoreOpenChangesStep) CreateAbortStep() Step {
  return nil
}

func (step RestoreOpenChangesStep) CreateContinueStep() Step {
  return nil
}

func (step RestoreOpenChangesStep) CreateUndoStep() Step {
  return nil
}

func (step RestoreOpenChangesStep) Run() error {
  return script.RunCommand([]string{"git", "stash", "pop"})
}
