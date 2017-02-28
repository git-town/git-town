package step

import (
  "github.com/Originate/gt/cmd/git"
  "github.com/Originate/gt/cmd/script"
)

type PushBranchStep struct {
  BranchName string
  Force bool
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
  if git.ShouldBranchBePushed(step.BranchName) {
    if step.Force {
      return script.RunCommand([]string{"git", "push", "-f", "origin", step.BranchName})
    } else if git.GetCurrentBranchName() == step.BranchName {
      return script.RunCommand([]string{"git", "push"})
    } else {
      return script.RunCommand([]string{"git", "push", "origin", step.BranchName})
    }
  }
  return nil
}
