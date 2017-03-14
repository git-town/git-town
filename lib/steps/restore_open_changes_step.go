package steps

import (
  "github.com/Originate/gt/lib/script"
)


type RestoreOpenChangesStep struct {}


func (step RestoreOpenChangesStep) CreateAbortStep() Step {
  return NoOpStep{}
}


func (step RestoreOpenChangesStep) CreateContinueStep() Step {
  return NoOpStep{}
}


func (step RestoreOpenChangesStep) CreateUndoStep() Step {
  return NoOpStep{}
}


func (step RestoreOpenChangesStep) Run() error {
  return script.RunCommand("git", "stash", "pop")
}
