package steps

import (
  "github.com/Originate/gt/cmd/git"
)


type RebaseTrackingBranchStep struct {}


func (step RebaseTrackingBranchStep) CreateAbortStep() Step {
  return AbortRebaseBranchStep{}
}


func (step RebaseTrackingBranchStep) CreateContinueStep() Step {
  return ContinueRebaseBranchStep{}
}


func (step RebaseTrackingBranchStep) CreateUndoStep() Step {
  return ResetToShaStep{Hard: true, Sha: git.GetCurrentSha()}
}


func (step RebaseTrackingBranchStep) Run() error {
  branchName := git.GetCurrentBranchName()
  if git.HasTrackingBranch(branchName) {
    step := RebaseBranchStep{BranchName: git.GetTrackingBranchName(branchName)}
    return step.Run()
  }
  return nil
}
