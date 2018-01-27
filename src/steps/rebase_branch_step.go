package steps

import (
	"github.com/Originate/git-town/src/git"
	"github.com/Originate/git-town/src/script"
)

// RebaseBranchStep rebases the current branch
// against the branch with the given name.
type RebaseBranchStep struct {
	NoOpStep
	BranchName string

	previousSha string
}

// CreateAbortStep returns the abort step for this step.
func (step *RebaseBranchStep) CreateAbortStep() Step {
	return &AbortRebaseBranchStep{}
}

// CreateContinueStep returns the continue step for this step.
func (step *RebaseBranchStep) CreateContinueStep() Step {
	return &ContinueRebaseBranchStep{}
}

// AddUndoSteps adds the undo steps for this step to the undo step list
func (step *RebaseBranchStep) AddUndoSteps(stepList *StepList) {
	stepList.Prepend(&ResetToShaStep{Hard: true, Sha: step.previousSha})
}

// Run executes this step.
func (step *RebaseBranchStep) Run() error {
	step.previousSha = git.GetCurrentSha()
	err := script.RunCommand("git", "rebase", step.BranchName)
	if err != nil {
		git.ClearCurrentBranchCache()
	}
	return err
}
