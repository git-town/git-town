package steps

import (
	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/hosting"
)

// RebaseBranchStep rebases the current branch
// against the branch with the given name.
type RebaseBranchStep struct {
	NoOpStep
	Branch      string
	previousSha string
}

func (step *RebaseBranchStep) CreateAbortStep() Step { //nolint:ireturn
	return &AbortRebaseBranchStep{}
}

func (step *RebaseBranchStep) CreateContinueStep() Step { //nolint:ireturn
	return &ContinueRebaseBranchStep{}
}

func (step *RebaseBranchStep) CreateUndoStep(repo *git.ProdRepo) (Step, error) { //nolint:ireturn
	return &ResetToShaStep{Hard: true, Sha: step.previousSha}, nil
}

func (step *RebaseBranchStep) Run(repo *git.ProdRepo, driver hosting.Driver) error {
	var err error
	step.previousSha, err = repo.Silent.CurrentSha()
	if err != nil {
		return err
	}
	err = repo.Logging.Rebase(step.Branch)
	if err != nil {
		repo.Silent.CurrentBranchCache.Invalidate()
	}
	return err
}
