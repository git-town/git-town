package step

import (
  "github.com/Originate/gt/cmd/git"
  "github.com/Originate/gt/cmd/script"
)

type ContinueRebaseBranchStep int

func (step ContinueRebaseBranchStep) CreateAbortStep() Step {
  return new(NoOpStep)
}

func (step ContinueRebaseBranchStep) CreateContinueStep() Step {
  return new(NoOpStep)
}

func (step ContinueRebaseBranchStep) CreateUndoStep() Step {
  return new(NoOpStep)
}

func (step ContinueRebaseBranchStep) Run() error {
  if git.IsRebaseInProgress() {
    return script.RunCommand([]string{"git", "rebase", "--continue"})
  }
  return nil
}
