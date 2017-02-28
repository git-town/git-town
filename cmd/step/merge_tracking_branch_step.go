package step

import (
  "github.com/Originate/gt/cmd/git"
)

type MergeTrackingBranchStep struct {}

func (step MergeTrackingBranchStep) CreateAbortStep() Step {
  return new(AbortMergeBranchStep)
}

func (step MergeTrackingBranchStep) CreateContinueStep() Step {
  return new(ContinueMergeBranchStep)
}

func (step MergeTrackingBranchStep) CreateUndoStep() Step {
  return ResetToShaStep{hard: true, sha: git.GetCurrentSha()}
}

func (step MergeTrackingBranchStep) Run() error {
  branchName := git.GetCurrentBranchName()
  if git.HasTrackingBranch(branchName) {
    step := MergeBranchStep{branchName: git.GetTrackingBranchName(branchName)}
    return step.Run()
  }
  return nil
}
