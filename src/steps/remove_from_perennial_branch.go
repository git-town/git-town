package steps

import "github.com/git-town/git-town/src/git"

// RemoveFromPerennialBranches removes the branch with the given name as a perennial branch
type RemoveFromPerennialBranches struct {
	NoOpStep
	BranchName string
}

// AddUndoSteps adds the undo steps for this step to the undo step list
func (step *RemoveFromPerennialBranches) AddUndoSteps(stepList *StepList) {
	stepList.Prepend(&AddToPerennialBranches{BranchName: step.BranchName})
}

// Run executes this step.
func (step *RemoveFromPerennialBranches) Run() error {
	git.Config().RemoveFromPerennialBranches(step.BranchName)
	return nil
}
