package steps

import (
	"github.com/git-town/git-town/v7/src/config"
	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/hosting"
)

// CreateTrackingBranchStep pushes the current branch up to origin
// and marks it as tracking the current branch.
type CreateTrackingBranchStep struct {
	EmptyStep
	Branch     string
	NoPushHook bool
}

func (step *CreateTrackingBranchStep) CreateUndoStep(repo *git.ProdRepo) (Step, error) {
	return &DeleteOriginBranchStep{Branch: step.Branch}, nil
}

func (step *CreateTrackingBranchStep) Run(repo *git.ProdRepo, connector hosting.Connector) error {
	return repo.Logging.PushBranch(git.PushArgs{
		Branch:     step.Branch,
		NoPushHook: step.NoPushHook,
		Remote:     config.OriginRemote,
	})
}
