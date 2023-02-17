package steps

import (
	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/hosting"
)

// CreateBranchStep creates a new branch
// but leaves the current branch unchanged.
type CreateBranchStep struct {
	NoOpStep
	Branch        string
	StartingPoint string
}

func (step *CreateBranchStep) CreateUndoStep(repo *git.ProdRepo) (Step, error) { //nolint:ireturn
	return &DeleteLocalBranchStep{Branch: step.Branch, Force: true}, nil
}

func (step *CreateBranchStep) Run(repo *git.ProdRepo, driver hosting.Driver) error {
	return repo.Logging.CreateBranch(step.Branch, step.StartingPoint)
}
