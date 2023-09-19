package steps

import (
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/git"
	"github.com/git-town/git-town/v9/src/hosting"
)

// ResetRemoteBranchToSHAStep sets the given remote branch to the given SHA.
type ResetRemoteBranchToSHAStep struct {
	Branch           domain.RemoteBranchName
	SHAToPush        domain.SHA // the SHA to reset the remote branch to
	SHAThatMustExist domain.SHA // the SHA that the branch must have at the remote in order to perform the reset without losing data
	EmptyStep
}

func (step *ResetRemoteBranchToSHAStep) Run(run *git.ProdRunner, _ hosting.Connector) error {
	currentRemoteSHA, err := run.Backend.SHAForBranch(step.Branch.BranchName())
	if err != nil {
		return err
	}
	if currentRemoteSHA == step.SHAToPush {
		return nil
	}
	return run.Frontend.ResetRemoteBranchToSHA(step.Branch, step.SHAToPush, step.SHAThatMustExist)
}
