package steps

import (
	"github.com/git-town/git-town/v9/src/config"
	"github.com/git-town/git-town/v9/src/git"
	"github.com/git-town/git-town/v9/src/hosting"
)

// CreateTrackingBranchStep pushes the current branch up to origin
// and marks it as tracking the current branch.
type CreateTrackingBranchStep struct {
	EmptyStep
	Branch     string
	NoPushHook bool
}

func (step *CreateTrackingBranchStep) CreateUndoStep(backend *git.BackendCommands) (Step, error) {
	return &DeleteOriginBranchStep{Branch: step.Branch}, nil
}

func (step *CreateTrackingBranchStep) Run(run *git.ProdRunner, connector hosting.Connector) error {
	return run.Frontend.PushBranch(git.PushArgs{
		Branch:     step.Branch,
		NoPushHook: step.NoPushHook,
		Remote:     config.OriginRemote,
	})
}
