package steps

import "github.com/Originate/git-town/lib/script"

// CreateBranchStep creates a new branch
// but leaves the current branch unchanged.
type CreateBranchStep struct {
	NoOpStep
	BranchName    string
	StartingPoint string
}

// CreateUndoStepBeforeRun returns the undo step for this step before it is run.
func (step CreateBranchStep) CreateUndoStepBeforeRun() Step {
	return DeleteLocalBranchStep{BranchName: step.BranchName}
}

// Run executes this step.
func (step CreateBranchStep) Run() error {
	return script.RunCommand("git", "branch", step.BranchName, step.StartingPoint)
}
