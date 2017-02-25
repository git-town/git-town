package step

import (
  "github.com/Originate/gt/cmd/script"
)

type AbortMergeBranchStep int

func (step AbortMergeBranchStep) CreateAbortStep() Step {
  return new(NoOpStep)
}

func (step AbortMergeBranchStep) CreateContinueStep() Step {
  return new(NoOpStep)
}

func (step AbortMergeBranchStep) CreateUndoStep() Step {
  return new(NoOpStep)
}

func (step AbortMergeBranchStep) Run() error {
  return script.RunCommand([]string{"git", "merge", "--abort"})
}
