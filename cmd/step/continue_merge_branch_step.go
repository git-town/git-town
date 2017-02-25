package step

import (
  "github.com/Originate/gt/cmd/git"
  "github.com/Originate/gt/cmd/script"
)

type ContinueMergeBranchStep int

func (step ContinueMergeBranchStep) CreateAbortStep() Step {
  return new(NoOpStep)
}

func (step ContinueMergeBranchStep) CreateContinueStep() Step {
  return new(NoOpStep)
}

func (step ContinueMergeBranchStep) CreateUndoStep() Step {
  return new(NoOpStep)
}

func (step ContinueMergeBranchStep) Run() error {
  if git.IsMergeInProgress() {
    return script.RunCommand([]string{"git", "commit", "--no-edit"})
  }
  return nil
}
