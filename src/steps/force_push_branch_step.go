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
	OriginalRemoteSHA domain.SHA // the SHA that the remote branch had before Git Town ran
	shaAfterPush      domain.SHA
	EmptyStep
}

func (step *ForcePushBranchStep) CreateUndoSteps(_ *git.BackendCommands) ([]Step, error) {
	return []Step{&ResetRemoteBranchToSHAStep{
		Branch:           step.Branch.RemoteName(),
		SHAToPush:        step.OriginalRemoteSHA,
		SHAThatMustExist: step.shaAfterPush,
	}}, nil
}

func (step *ForcePushBranchStep) Run(run *git.ProdRunner, _ hosting.Connector) error {
	shouldPush, err := run.Backend.ShouldPushBranch(step.Branch, step.Branch.RemoteName())
	if err != nil {
		return err
	}
	if !shouldPush && !run.Config.DryRun {
		return nil
	}
	step.shaAfterPush, err = run.Backend.SHAForBranch(step.Branch.BranchName())
	if err != nil {
		return err
	}
	return run.Frontend.ForcePushBranch(step.NoPushHook)
}
