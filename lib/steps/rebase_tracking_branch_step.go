package steps

import (
  "github.com/Originate/gt/lib/git"
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
    return RebaseBranchStep{BranchName: git.GetTrackingBranchName(branchName)}.Run()
  }
  return nil
}
