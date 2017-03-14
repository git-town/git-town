package steps

import (
  "github.com/Originate/gt/lib/git"
  "github.com/Originate/gt/lib/script"
)


type ContinueMergeBranchStep struct {}


func (step ContinueMergeBranchStep) CreateAbortStep() Step {
  return NoOpStep{}
}


func (step ContinueMergeBranchStep) CreateContinueStep() Step {
  return step
}


func (step ContinueMergeBranchStep) CreateUndoStep() Step {
  return NoOpStep{}
}


func (step ContinueMergeBranchStep) Run() error {
  if git.IsMergeInProgress() {
    return script.RunCommand("git", "commit", "--no-edit")
  }
  return nil
}
