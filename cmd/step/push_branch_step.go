package step

import (
  "github.com/Originate/gt/cmd/git"
  "github.com/Originate/gt/cmd/script"
)

type PushBranchStep struct {
  branchName string
  force bool
}

func (step PushBranchStep) CreateAbortStep() Step {
  return nil
}

func (step PushBranchStep) CreateContinueStep() Step {
  return nil
}

func (step PushBranchStep) CreateUndoStep() Step {
  return nil
}

func (step PushBranchStep) Run() error {
  if git.ShouldBranchBePushed(step.branchName) {
    if step.force {
      return script.RunCommand([]string{"git", "push", "-f", "origin", step.branchName})
    } else if git.GetCurrentBranchName() == step.branchName {
      return script.RunCommand([]string{"git", "push"})
    } else {
      return script.RunCommand([]string{"git", "push", "origin", step.branchName})
    }
  }
  return nil
}
