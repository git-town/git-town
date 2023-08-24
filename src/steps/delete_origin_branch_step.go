package steps

import (
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/git"
	"github.com/git-town/git-town/v9/src/hosting"
)

// DeleteOriginBranchStep deletes the current branch from the origin remote.
type DeleteOriginBranchStep struct {
	Branch     domain.LocalBranchName
	NoPushHook bool
	branchSha  domain.SHA `exhaustruct:"optional"`
	EmptyStep
}

func (step *DeleteOriginBranchStep) CreateUndoSteps(_ *git.BackendCommands) ([]Step, error) {
	return []Step{&CreateRemoteBranchStep{Branch: step.Branch, Sha: step.branchSha, NoPushHook: step.NoPushHook}}, nil
}

func (step *DeleteOriginBranchStep) Run(run *git.ProdRunner, _ hosting.Connector) error {
	remoteBranch := step.Branch.AtRemote(domain.OriginRemote) // TODO: inject git.Branches somehow and look the name of the actual tracking brach in it
	var err error
	step.branchSha, err = run.Backend.ShaForBranch(remoteBranch.BranchName())
	if err != nil {
		return err
	}
	return run.Frontend.DeleteRemoteBranch(step.Branch)
}
