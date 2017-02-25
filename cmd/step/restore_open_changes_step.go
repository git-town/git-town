package step

import (
  "github.com/Originate/gt/cmd/script"
)

type RestoreOpenChangesStep int

func (step RestoreOpenChangesStep) CreateAbortStep() Step {
  return new(NoOpStep)
}

func (step RestoreOpenChangesStep) CreateContinueStep() Step {
  return new(NoOpStep)
}

func (step RestoreOpenChangesStep) CreateUndoStep() Step {
  return new(NoOpStep)
}

func (step RestoreOpenChangesStep) Run() error {
  return script.RunCommand([]string{"git", "stash", "pop"})
}
