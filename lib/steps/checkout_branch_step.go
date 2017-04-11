package steps

import (
	"github.com/Originate/git-town/lib/git"
	"github.com/Originate/git-town/lib/script"
)

type CheckoutBranchStep struct {
	NoAutomaticAbortOnError
	NoUndoStepAfterRun
	BranchName string
}

func (step CheckoutBranchStep) CreateAbortStep() Step {
	return NoOpStep{}
}

func (step CheckoutBranchStep) CreateContinueStep() Step {
	return NoOpStep{}
}

func (step CheckoutBranchStep) CreateUndoStepBeforeRun() Step {
	return CheckoutBranchStep{BranchName: git.GetCurrentBranchName()}
}

func (step CheckoutBranchStep) Run() error {
	if git.GetCurrentBranchName() != step.BranchName {
		return script.RunCommand("git", "checkout", step.BranchName)
	}
	return nil
}
