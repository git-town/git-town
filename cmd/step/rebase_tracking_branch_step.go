package step

import (
  "github.com/Originate/gt/cmd/git"
)

type RebaseTrackingBranchStep struct {}

func (step RebaseTrackingBranchStep) CreateAbortStep() Step {
  return new(AbortRebaseBranchStep)
}

func (step RebaseTrackingBranchStep) CreateContinueStep() Step {
  return new(ContinueRebaseBranchStep)
}

func (step RebaseTrackingBranchStep) CreateUndoStep() Step {
  return ResetToShaStep{hard: true, sha: git.GetCurrentSha()}
}

func (step RebaseTrackingBranchStep) Run() error {
  branchName := git.GetCurrentBranchName()
  if git.HasTrackingBranch(branchName) {
    step := RebaseBranchStep{branchName: git.GetTrackingBranchName(branchName)}
    return step.Run()
  }
  return nil
}
