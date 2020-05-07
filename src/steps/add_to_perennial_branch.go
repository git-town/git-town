package steps

import "github.com/git-town/git-town/src/git"

// AddToPerennialBranches adds the branch with the given name as a perennial branch
type AddToPerennialBranches struct {
	NoOpStep
	BranchName string
}

// AddUndoSteps adds the undo steps for this step to the undo step list
func (step *AddToPerennialBranches) AddUndoSteps(stepList *StepList) {
	stepList.Prepend(&RemoveFromPerennialBranches{BranchName: step.BranchName})
}

// Run executes this step.
func (step *AddToPerennialBranches) Run() error {
	git.Config().AddToPerennialBranches(step.BranchName)
	return nil
}
