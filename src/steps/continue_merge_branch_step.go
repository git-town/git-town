//nolint:ireturn
package steps

import (
	"github.com/git-town/git-town/v7/src/drivers"
	"github.com/git-town/git-town/v7/src/git"
)

// ContinueMergeBranchStep finishes an ongoing merge conflict
// assuming all conflicts have been resolved by the user.
type ContinueMergeBranchStep struct {
	NoOpStep
}

// CreateAbortStep returns the abort step for this step.
func (step *ContinueMergeBranchStep) CreateAbortStep() Step {
	return &NoOpStep{}
}

// CreateContinueStep returns the continue step for this step.
func (step *ContinueMergeBranchStep) CreateContinueStep() Step {
	return step
}

// Run executes this step.
func (step *ContinueMergeBranchStep) Run(repo *git.ProdRepo, driver drivers.CodeHostingDriver) error {
	hasMergeInprogress, err := repo.Silent.HasMergeInProgress()
	if err != nil {
		return err
	}
	if hasMergeInprogress {
		return repo.Logging.CommitNoEdit()
	}
	return nil
}
