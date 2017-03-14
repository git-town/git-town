package steps

import (
  "github.com/Originate/git-town/lib/git"
  "github.com/Originate/git-town/lib/script"
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
    return script.RunCommand("git", "rebase", "--continue")
  }
  return nil
}
