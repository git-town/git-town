package step

import (
  "github.com/Originate/gt/cmd/script"
)

type CreateTrackingBranchStep struct {
  BranchName string
}

func (step CreateTrackingBranchStep) CreateAbortStep() Step {
  return NoOpStep{}
}

func (step CreateTrackingBranchStep) CreateContinueStep() Step {
  return NoOpStep{}
}

func (step CreateTrackingBranchStep) CreateUndoStep() Step {
  return NoOpStep{} // TODO delete remote branch
}

func (step CreateTrackingBranchStep) Run() error {
  return script.RunCommand([]string{"git", "push", "-u", "origin", step.BranchName})
}
