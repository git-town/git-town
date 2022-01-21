package steps

import (
	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/hosting"
)

// MergeBranchStep merges the branch with the given name into the current branch.
type MergeBranchStep struct {
	NoOpStep
	BranchName  string
	previousSha string
}

func (step *MergeBranchStep) CreateAbortStep() Step { //nolint:ireturn
	return &AbortMergeBranchStep{}
}

func (step *MergeBranchStep) CreateContinueStep() Step { //nolint:ireturn
	return &ContinueMergeBranchStep{}
}

func (step *MergeBranchStep) CreateUndoStep(repo *git.ProdRepo) (Step, error) { //nolint:ireturn
	return &ResetToShaStep{Hard: true, Sha: step.previousSha}, nil
}

func (step *MergeBranchStep) Run(repo *git.ProdRepo, driver hosting.Driver) (err error) {
	step.previousSha, err = repo.Silent.CurrentSha()
	if err != nil {
		return err
	}
	return repo.Logging.MergeBranchNoEdit(step.BranchName)
}
