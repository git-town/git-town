package steps

import (
	"github.com/Originate/git-town/lib/git"
	"github.com/Originate/git-town/lib/script"
)

// CreateAndCheckoutBranchStep creates a new branch and makes it the current one.
type CreateAndCheckoutBranchStep struct {
	NoOpStep
	BranchName       string
	ParentBranchName string
}

// Run executes this step.
func (step CreateAndCheckoutBranchStep) Run() error {
	git.SetParentBranch(step.BranchName, step.ParentBranchName)
	return script.RunCommand("git", "checkout", "-b", step.BranchName, step.ParentBranchName)
}
