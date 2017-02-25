package step

import (
  "github.com/Originate/gt/cmd/script"
)

type AbortMergeBranchStep int

func (step AbortMergeBranchStep) CreateAbortStep() Step {
  return nil
}

func (step AbortMergeBranchStep) CreateContinueStep() Step {
  return nil
}

func (step AbortMergeBranchStep) CreateUndoStep() Step {
  return nil
}

func (step AbortMergeBranchStep) Run() error {
  return script.RunCommand([]string{"git", "merge", "--abort"})
}
