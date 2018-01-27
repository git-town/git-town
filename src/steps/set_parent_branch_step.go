package steps

import "github.com/Originate/git-town/src/git"

// SetParentBranchStep registers the branch with the given name as a parent
// of the branch with the other given name.
type SetParentBranchStep struct {
	NoOpStep
	BranchName       string
	ParentBranchName string

	previousParent string
}

// AddUndoSteps adds the undo steps for this step to the undo step list
func (step *SetParentBranchStep) AddUndoSteps(stepList *StepList) {
	if step.previousParent == "" {
		stepList.Prepend(&DeleteParentBranchStep{BranchName: step.BranchName})
	} else {
		stepList.Prepend(&SetParentBranchStep{BranchName: step.BranchName, ParentBranchName: step.previousParent})
	}
}

// Run executes this step.
func (step *SetParentBranchStep) Run() error {
	step.previousParent = git.GetParentBranch(step.BranchName)
	git.SetParentBranch(step.BranchName, step.ParentBranchName)
	return nil
}
