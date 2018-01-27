package steps

import "github.com/Originate/git-town/src/git"

// DeleteParentBranchStep removes the parent branch entry in the Git Town configuration.
type DeleteParentBranchStep struct {
	NoOpStep
	BranchName string

	previousParent string
}

// AddUndoSteps adds the undo steps for this step to the undo step list
func (step *DeleteParentBranchStep) AddUndoSteps(stepList *StepList) {
	if step.previousParent == "" {
		return
	}
	stepList.Prepend(&SetParentBranchStep{BranchName: step.BranchName, ParentBranchName: step.previousParent})
}

// Run executes this step.
func (step *DeleteParentBranchStep) Run() error {
	step.previousParent = git.GetParentBranch(step.BranchName)
	git.DeleteParentBranch(step.BranchName)
	return nil
}
