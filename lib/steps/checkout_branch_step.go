package steps

import (
  "github.com/Originate/gt/lib/git"
  "github.com/Originate/gt/lib/script"
)


type CheckoutBranchStep struct {
  BranchName string
}


func (step CheckoutBranchStep) CreateAbortStep() Step {
  return NoOpStep{}
}


func (step CheckoutBranchStep) CreateContinueStep() Step {
  return NoOpStep{}
}


func (step CheckoutBranchStep) CreateUndoStep() Step {
  return CheckoutBranchStep{BranchName: git.GetCurrentBranchName()}
}


func (step CheckoutBranchStep) Run() error {
  if git.GetCurrentBranchName() != step.BranchName {
    return script.RunCommand("git", "checkout", step.BranchName)
  }
  return nil
}
