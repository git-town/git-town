package step

import (
  "github.com/Originate/gt/cmd/script"
)

type CreateTrackingBranchStep struct {
  branchName string
}

func (step CreateTrackingBranchStep) CreateAbortStep() Step {
  return new(NoOpStep)
}

func (step CreateTrackingBranchStep) CreateContinueStep() Step {
  return new(NoOpStep)
}

func (step CreateTrackingBranchStep) CreateUndoStep() Step {
  return new(NoOpStep) // TODO delete remote branch
}

func (step CreateTrackingBranchStep) Run() error {
  return script.RunCommand([]string{"git", "push", "origin", step.branchName})
}
