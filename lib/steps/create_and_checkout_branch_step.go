package steps

import (
	"github.com/Originate/git-town/lib/config"
	"github.com/Originate/git-town/lib/script"
)

type CreateAndCheckoutBranchStep struct {
	BranchName       string
	ParentBranchName string
}

func (step CreateAndCheckoutBranchStep) CreateAbortStep() Step {
	return NoOpStep{}
}

func (step CreateAndCheckoutBranchStep) CreateContinueStep() Step {
	return NoOpStep{}
}

func (step CreateAndCheckoutBranchStep) CreateUndoStep() Step {
	return NoOpStep{}
}

func (step CreateAndCheckoutBranchStep) GetAutomaticAbortErrorMessage() string {
	return ""
}

func (step CreateAndCheckoutBranchStep) Run() error {
	config.SetParentBranch(step.BranchName, step.ParentBranchName)
	return script.RunCommand("git", "checkout", "-b", step.BranchName, step.ParentBranchName)
}

func (step CreateAndCheckoutBranchStep) ShouldAutomaticallyAbortOnError() bool {
	return false
}
