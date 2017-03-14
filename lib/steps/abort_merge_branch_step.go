package steps

import (
  "github.com/Originate/git-town/lib/script"
)


type AbortMergeBranchStep struct {}


func (step AbortMergeBranchStep) CreateAbortStep() Step {
  return NoOpStep{}
}


func (step AbortMergeBranchStep) CreateContinueStep() Step {
  return NoOpStep{}
}


func (step AbortMergeBranchStep) CreateUndoStep() Step {
  return NoOpStep{}
}


func (step AbortMergeBranchStep) Run() error {
  return script.RunCommand("git", "merge", "--abort")
}
