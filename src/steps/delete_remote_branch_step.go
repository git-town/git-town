package steps

import (
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/git"
	"github.com/git-town/git-town/v9/src/hosting"
)

// DeleteRemoteBranchStep deletes the given branch from the origin remote.
// TODO: split this step type up into two and delete IsTracking.
type DeleteRemoteBranchStep struct {
	Branch     domain.LocalBranchName
	Remote     domain.Remote
	IsTracking bool
	NoPushHook bool
	branchSha  domain.SHA `exhaustruct:"optional"`
	EmptyStep
}

func (step *DeleteRemoteBranchStep) CreateUndoSteps(_ *git.BackendCommands) ([]Step, error) {
	if step.IsTracking {
		return []Step{&CreateTrackingBranchStep{Branch: step.Branch, NoPushHook: step.NoPushHook}}, nil
	}
	return []Step{&CreateRemoteBranchStep{Branch: step.Branch, Sha: step.branchSha, NoPushHook: step.NoPushHook}}, nil
}

func (step *DeleteRemoteBranchStep) Run(run *git.ProdRunner, _ hosting.Connector) error {
	if !step.IsTracking {
		trackingBranch := step.Branch.AtRemote(domain.OriginRemote) // TODO: inject git.Branches somehow and look the name of the actual tracking brach in it
		var err error
		step.branchSha, err = run.Backend.ShaForBranch(trackingBranch.BranchName())
		if err != nil {
			return err
		}
	}
	return run.Frontend.DeleteRemoteBranch(step.Branch)
}
