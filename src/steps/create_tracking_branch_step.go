package steps

import (
	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/hosting"
)

// CreateTrackingBranchStep pushes the current branch up to origin
// and marks it as tracking the current branch.
type CreateTrackingBranchStep struct {
	NoOpStep
	BranchName   string
	NoPushVerify bool
}

func (step *CreateTrackingBranchStep) CreateUndoStep(repo *git.ProdRepo) (Step, error) { //nolint:ireturn
	return &DeleteOriginBranchStep{BranchName: step.BranchName}, nil
}

func (step *CreateTrackingBranchStep) Run(repo *git.ProdRepo, driver hosting.Driver) error {
	return repo.Logging.PushBranch(git.PushBranchArgs{
		BranchName:   step.BranchName,
		NoPushVerify: step.NoPushVerify,
		ToOrigin:     true,
	})
}
