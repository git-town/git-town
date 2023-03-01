package steps

import (
	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/hosting"
)

// MergeBranchStep merges the branch with the given name into the current branch.
type MergeBranchStep struct {
	NoOpStep
	Branch      string
	previousSha string
}

func (step *MergeBranchStep) CreateAbortStep() Step {
	return &AbortMergeStep{}
}

func (step *MergeBranchStep) CreateContinueStep() Step {
	return &ContinueMergeStep{}
}

func (step *MergeBranchStep) CreateUndoStep(repo *git.ProdRepo) (Step, error) {
	return &ResetToShaStep{Hard: true, Sha: step.previousSha}, nil
}

func (step *MergeBranchStep) Run(repo *git.ProdRepo, connector hosting.Connector) error {
	var err error
	step.previousSha, err = repo.Silent.CurrentSha()
	if err != nil {
		return err
	}
	return repo.Logging.MergeBranchNoEdit(step.Branch)
}
