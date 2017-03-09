package steps

import (
  "github.com/Originate/gt/cmd/git"
  "github.com/Originate/gt/cmd/script"
)


type RebaseBranchStep struct {
  BranchName string
}


func (step RebaseBranchStep) CreateAbortStep() Step {
  return AbortRebaseBranchStep{}
}


func (step RebaseBranchStep) CreateContinueStep() Step {
  return ContinueRebaseBranchStep{}
}


func (step RebaseBranchStep) CreateUndoStep() Step {
  return ResetToShaStep{Hard: true, Sha: git.GetCurrentSha()}
}


func (step RebaseBranchStep) Run() error {
  return script.RunCommand([]string{"git", "rebase", step.BranchName})
}
