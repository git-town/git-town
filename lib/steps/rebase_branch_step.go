package steps

import (
	"github.com/Originate/git-town/lib/git"
	"github.com/Originate/git-town/lib/script"
)

type RebaseBranchStep struct {
	NoAutomaticAbortOnError
	NoUndoStepAfterRun
	BranchName string
}

func (step RebaseBranchStep) CreateAbortStep() Step {
	return AbortRebaseBranchStep{}
}

func (step RebaseBranchStep) CreateContinueStep() Step {
	return ContinueRebaseBranchStep{}
}

func (step RebaseBranchStep) CreateUndoStepBeforeRun() Step {
	return ResetToShaStep{Hard: true, Sha: git.GetCurrentSha()}
}

func (step RebaseBranchStep) Run() error {
	return script.RunCommand("git", "rebase", step.BranchName)
}
