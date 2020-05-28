package steps

import (
	"github.com/git-town/git-town/src/git"
)

// MergeBranchStep merges the branch with the given name into the current branch
type MergeBranchStep struct {
	NoOpStep
	BranchName string

	previousSha string
}

// CreateAbortStep returns the abort step for this step.
func (step *MergeBranchStep) CreateAbortStep() Step {
	return &AbortMergeBranchStep{}
}

// CreateContinueStep returns the continue step for this step.
func (step *MergeBranchStep) CreateContinueStep() Step {
	return &ContinueMergeBranchStep{}
}

// CreateUndoStep returns the undo step for this step.
func (step *MergeBranchStep) CreateUndoStep() Step {
	return &ResetToShaStep{Hard: true, Sha: step.previousSha}
}

// Run executes this step.
func (step *MergeBranchStep) Run(repo *git.ProdRepo) (err error) {
	step.previousSha, err = repo.Silent.CurrentSha()
	if err != nil {
		return err
	}
	return repo.Logging.MergeBranchNoEdit(step.BranchName)
}
