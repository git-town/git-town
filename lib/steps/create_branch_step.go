package steps

import "github.com/Originate/git-town/lib/script"

// CreateBranchStep creates a new branch
// but leaves the current branch unchanged.
type CreateBranchStep struct {
	NoExpectedError
	NoUndoStep
	BranchName    string
	StartingPoint string
}

// Run executes this step.
func (step CreateBranchStep) Run() error {
	return script.RunCommand("git", "branch", step.BranchName, step.StartingPoint)
}
