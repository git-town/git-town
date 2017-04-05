package steps

import (
	"github.com/Originate/git-town/lib/git"
	"github.com/Originate/git-town/lib/script"
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

func (step MergeBranchStep) GetAutomaticAbortErrorMessage() string {
	return ""
}

func (step MergeBranchStep) Run() error {
	return script.RunCommand("git", "merge", "--no-edit", step.BranchName)
}

func (step MergeBranchStep) ShouldAutomaticallyAbortOnError() bool {
	return false
}
