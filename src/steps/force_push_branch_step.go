package steps

import (
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/git"
	"github.com/git-town/git-town/v9/src/hosting"
)

// ForcePushBranchStep force-pushes the branch with the given name to the origin remote.
type ForcePushBranchStep struct {
	Branch            domain.LocalBranchName
	NoPushHook        bool
	originalRemoteSHA domain.SHA // the SHA that the remote branch had before Git Town ran
	shaAfterPush      domain.SHA
	EmptyStep
}

func (step *ForcePushBranchStep) CreateUndoSteps(_ *git.BackendCommands) ([]Step, error) {
	return []Step{&ResetRemoteBranchToSHAStep{
		Branch:           step.Branch.RemoteBranch(),
		SHAToPush:        step.originalRemoteSHA,
		SHAThatMustExist: step.shaAfterPush,
	}}, nil
}

func (step *ForcePushBranchStep) Run(run *git.ProdRunner, _ hosting.Connector) error {
	remoteBranch := step.Branch.RemoteBranch()
	shouldPush, err := run.Backend.ShouldPushBranch(step.Branch, remoteBranch)
	if err != nil {
		return err
	}
	if !shouldPush && !run.Config.DryRun {
		return nil
	}
	step.originalRemoteSHA, err = run.Backend.SHAForBranch(remoteBranch.BranchName())
	if err != nil {
		return err
	}
	step.shaAfterPush, err = run.Backend.SHAForBranch(step.Branch.BranchName())
	if err != nil {
		return err
	}
	return run.Frontend.ForcePushBranch(step.NoPushHook)
}
