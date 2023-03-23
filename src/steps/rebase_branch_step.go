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

func (step *RebaseBranchStep) CreateUndoStep(backend *git.BackendCommands) (Step, error) {
	return &ResetToShaStep{Hard: true, Sha: step.previousSha}, nil
}

func (step *RebaseBranchStep) Run(repo *git.ProdRepo, connector hosting.Connector) error {
	var err error
	step.previousSha, err = repo.Backend.CurrentSha()
	if err != nil {
		return err
	}
	err = repo.Frontend.Rebase(step.Branch)
	if err != nil {
		repo.Config.CurrentBranchCache.Invalidate()
	}
	return err
}
