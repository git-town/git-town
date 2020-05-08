package steps

import "github.com/git-town/git-town/src/script"

// CreateBranchStep creates a new branch
// but leaves the current branch unchanged.
type CreateBranchStep struct {
	NoOpStep
	BranchName    string
	StartingPoint string
}

// CreateUndoStep returns the undo step for this step.
func (step *CreateBranchStep) CreateUndoStep() Step {
	return &DeleteLocalBranchStep{BranchName: step.BranchName}
}

// Run executes this step.
func (step *CreateBranchStep) Run() error {
	return script.RunCommand("git", "branch", step.BranchName, step.StartingPoint)
}
