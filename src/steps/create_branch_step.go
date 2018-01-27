package steps

import "github.com/Originate/git-town/src/script"

// CreateBranchStep creates a new branch
// but leaves the current branch unchanged.
type CreateBranchStep struct {
	NoOpStep
	BranchName    string
	StartingPoint string
}

// AddUndoSteps adds the undo steps for this step to the undo step list
func (step *CreateBranchStep) AddUndoSteps(stepList *StepList) {
	stepList.Prepend(&DeleteLocalBranchStep{BranchName: step.BranchName})
}

// Run executes this step.
func (step *CreateBranchStep) Run() error {
	return script.RunCommand("git", "branch", step.BranchName, step.StartingPoint)
}
