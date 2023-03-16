package steps

import (
	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/hosting"
)

// RebaseBranchStep rebases the current branch
// against the branch with the given name.
type RebaseBranchStep struct {
	EmptyStep
	Branch      string
	previousSha string
}

func (step *RebaseBranchStep) CreateAbortStep() Step {
	return &AbortRebaseStep{}
}

func (step *RebaseBranchStep) CreateContinueStep() Step {
	return &ContinueRebaseStep{}
}

func (step *RebaseBranchStep) CreateUndoStep(repo *git.PublicRepo) (Step, error) {
	return &ResetToShaStep{Hard: true, Sha: step.previousSha}, nil
}

func (step *RebaseBranchStep) Run(repo *git.PublicRepo, connector hosting.Connector) error {
	var err error
	step.previousSha, err = repo.CurrentSha()
	if err != nil {
		return err
	}
	err = repo.Rebase(step.Branch)
	if err != nil {
		repo.CurrentBranchCache.Invalidate()
	}
	return err
}
