package step

import (
  "github.com/Originate/gt/cmd/config"
  "github.com/Originate/gt/cmd/script"
)

type CreateAndCheckoutBranchStep struct {
  BranchName string
  ParentBranchName string
}

func (step CreateAndCheckoutBranchStep) CreateAbortStep() Step {
  return nil
}

func (step CreateAndCheckoutBranchStep) CreateContinueStep() Step {
  return nil
}

func (step CreateAndCheckoutBranchStep) CreateUndoStep() Step {
  return nil
}

func (step CreateAndCheckoutBranchStep) Run() error {
  config.StoreParentBranch(step.BranchName, step.ParentBranchName)
  return script.RunCommand([]string{"git", "checkout", "-b", step.BranchName, step.ParentBranchName})
}
