package step

import (
  "github.com/Originate/gt/cmd/git"
  "github.com/Originate/gt/cmd/script"
)

type ContinueMergeBranchStep struct {}

func (step ContinueMergeBranchStep) CreateAbortStep() Step {
  return nil
}

func (step ContinueMergeBranchStep) CreateContinueStep() Step {
  return nil
}

func (step ContinueMergeBranchStep) CreateUndoStep() Step {
  return nil
}

func (step ContinueMergeBranchStep) Run() error {
  if git.IsMergeInProgress() {
    return script.RunCommand([]string{"git", "commit", "--no-edit"})
  }
  return nil
}
