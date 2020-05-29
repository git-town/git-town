package steps

import (
	"github.com/git-town/git-town/src/git"
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

// CreateUndoStep returns the undo step for this step.
func (step *RebaseBranchStep) CreateUndoStep() Step {
	return &ResetToShaStep{Hard: true, Sha: step.previousSha}
}

// Run executes this step.
func (step *RebaseBranchStep) Run(repo *git.ProdRepo) error {
	step.previousSha = git.GetCurrentSha()
	err := repo.Logging.Rebase(step.BranchName)
	if err != nil {
		git.ClearCurrentBranchCache()
	}
	return err
}
