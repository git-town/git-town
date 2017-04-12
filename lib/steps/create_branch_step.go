package steps

import "github.com/Originate/git-town/lib/script"

// CreateBranchStep creates a new branch
// but leaves the current branch unchanged.
type CreateBranchStep struct {
	BranchName    string
	StartingPoint string
}

// CreateAbortStep returns the abort step for this step.
func (step CreateBranchStep) CreateAbortStep() Step {
	return NoOpStep{}
}

// CreateContinueStep returns the continue step for this step.
func (step CreateBranchStep) CreateContinueStep() Step {
	return NoOpStep{}
}

// CreateUndoStep returns the undo step for this step.
func (step CreateBranchStep) CreateUndoStep() Step {
	return NoOpStep{}
}

// Run executes this step.
func (step CreateBranchStep) Run() error {
	return script.RunCommand("git", "branch", step.BranchName, step.StartingPoint)
}
