package steps

import (
  "github.com/Originate/git-town/lib/script"
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
  return script.RunCommand("git", "rebase", "--abort")
}
