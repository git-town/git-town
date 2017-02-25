package step

import (
  "github.com/Originate/gt/cmd/git"
  "github.com/Originate/gt/cmd/script"
)

type RebaseBranchStep struct {
  branchName string
}

func (step RebaseBranchStep) CreateAbortStep() Step {
  return new(AbortRebaseBranchStep)
}

func (step RebaseBranchStep) CreateContinueStep() Step {
  return new(ContinueRebaseBranchStep)
}

func (step RebaseBranchStep) CreateUndoStep() Step {
  return ResetToShaStep{hard: true, sha: git.GetCurrentSha()}
}

func (step RebaseBranchStep) Run() error {
  return script.RunCommand([]string{"git", "rebase", "--no-edit", step.branchName})
}
