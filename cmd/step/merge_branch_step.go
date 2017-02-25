package step

import (
  "github.com/Originate/gt/cmd/git"
  "github.com/Originate/gt/cmd/script"
)

type MergeBranchStep struct {
  branchName string
}

func (step MergeBranchStep) CreateAbortStep() Step {
  return new(AbortMergeBranchStep)
}

func (step MergeBranchStep) CreateContinueStep() Step {
  return new(ContinueMergeBranchStep)
}

func (step MergeBranchStep) CreateUndoStep() Step {
  return ResetToShaStep{hard: true, sha: git.GetCurrentSha()}
}

func (step MergeBranchStep) Run() error {
  return script.RunCommand([]string{"git", "merge", "--no-edit", step.branchName})
}
