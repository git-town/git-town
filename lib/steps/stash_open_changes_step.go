package steps

import (
  "github.com/Originate/git-town/lib/script"
)


type StashOpenChangesStep struct {}


func (step StashOpenChangesStep) CreateAbortStep() Step {
  return NoOpStep{}
}


func (step StashOpenChangesStep) CreateContinueStep() Step {
  return NoOpStep{}
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
