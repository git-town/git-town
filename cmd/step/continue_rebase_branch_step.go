package step

import (
  "github.com/Originate/gt/cmd/git"
  "github.com/Originate/gt/cmd/script"
)

type ContinueRebaseBranchStep struct {}

func (step ContinueRebaseBranchStep) CreateAbortStep() Step {
  return nil
}

func (step ContinueRebaseBranchStep) CreateContinueStep() Step {
  return nil
}

func (step ContinueRebaseBranchStep) CreateUndoStep() Step {
  return nil
}

func (step ContinueRebaseBranchStep) Run() error {
  if git.IsRebaseInProgress() {
    return script.RunCommand([]string{"git", "rebase", "--continue"})
  }
  return nil
}
