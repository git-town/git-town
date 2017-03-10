package steps

import (
  "github.com/Originate/gt/lib/git"
  "github.com/Originate/gt/lib/script"
)


type ContinueRebaseBranchStep struct {}


func (step ContinueRebaseBranchStep) CreateAbortStep() Step {
  return NoOpStep{}
}


func (step ContinueRebaseBranchStep) CreateContinueStep() Step {
  return step
}


func (step ContinueRebaseBranchStep) CreateUndoStep() Step {
  return NoOpStep{}
}


func (step ContinueRebaseBranchStep) Run() error {
  if git.IsRebaseInProgress() {
    return script.RunCommand([]string{"git", "rebase", "--continue"})
  }
  return nil
}
