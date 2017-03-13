package steps

import (
  "github.com/Originate/gt/lib/git"
  "github.com/Originate/gt/lib/script"
)


type MergeBranchStep struct {
  BranchName string
}


func (step MergeBranchStep) CreateAbortStep() Step {
  return AbortMergeBranchStep{}
}


func (step MergeBranchStep) CreateContinueStep() Step {
  return ContinueMergeBranchStep{}
}


func (step MergeBranchStep) CreateUndoStep() Step {
  return ResetToShaStep{Hard: true, Sha: git.GetCurrentSha()}
}


func (step MergeBranchStep) Run() error {
  return script.RunCommand("git", "merge", "--no-edit", step.BranchName)
}
