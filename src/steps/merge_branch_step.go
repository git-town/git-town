//nolint:ireturn
package steps

import (
	"github.com/git-town/git-town/v7/src/drivers"
	"github.com/git-town/git-town/v7/src/git"
)

// MergeBranchStep merges the branch with the given name into the current branch.
type MergeBranchStep struct {
	NoOpStep
	BranchName string

	previousSha string
}

func (step *MergeBranchStep) CreateAbortStep() Step {
	return &AbortMergeBranchStep{}
}

func (step *MergeBranchStep) CreateContinueStep() Step {
	return &ContinueMergeBranchStep{}
}

func (step *MergeBranchStep) CreateUndoStep(repo *git.ProdRepo) (Step, error) {
	return &ResetToShaStep{Hard: true, Sha: step.previousSha}, nil
}

func (step *MergeBranchStep) Run(repo *git.ProdRepo, driver drivers.CodeHostingDriver) (err error) {
	step.previousSha, err = repo.Silent.CurrentSha()
	if err != nil {
		return err
	}
	return repo.Logging.MergeBranchNoEdit(step.BranchName)
}
