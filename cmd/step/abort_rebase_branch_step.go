package step

import (
  "github.com/Originate/gt/cmd/script"
)

type AbortRebaseBranchStep int

func (step AbortRebaseBranchStep) CreateAbortStep() Step {
  return new(NoOpStep)
}

func (step AbortRebaseBranchStep) CreateContinueStep() Step {
  return new(NoOpStep)
}

func (step AbortRebaseBranchStep) CreateUndoStep() Step {
  return new(NoOpStep)
}

func (step AbortRebaseBranchStep) Run() error {
  return script.RunCommand([]string{"git", "rebase", "--abort"})
}
