package step

import (
  "github.com/Originate/gt/cmd/git"
  "github.com/Originate/gt/cmd/script"
)

type CheckoutBranchStep struct {
  branchName string
}

func (step CheckoutBranchStep) CreateAbortStep() Step {
  return new(NoOpStep)
}

func (step CheckoutBranchStep) CreateContinueStep() Step {
  return new(NoOpStep)
}

func (step CheckoutBranchStep) CreateUndoStep() Step {
  return CheckoutBranchStep{branchName: git.GetCurrentBranchName()}
}

func (step CheckoutBranchStep) Run() error {
  return script.RunCommand([]string{"git", "checkout", step.branchName})
}
