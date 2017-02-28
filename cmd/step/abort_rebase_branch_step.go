package step

import (
  "github.com/Originate/gt/cmd/script"
)

type AbortRebaseBranchStep struct {}

func (step AbortRebaseBranchStep) CreateAbortStep() Step {
  return NoOpStep{}
}

func (step AbortRebaseBranchStep) CreateContinueStep() Step {
  return NoOpStep{}
}

func (step AbortRebaseBranchStep) CreateUndoStep() Step {
  return NoOpStep{}
}

func (step AbortRebaseBranchStep) Run() error {
  return script.RunCommand([]string{"git", "rebase", "--abort"})
}
