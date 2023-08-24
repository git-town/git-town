package steps

import (
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/git"
	"github.com/git-town/git-town/v9/src/hosting"
)

// PushCurrentBranchStep pushes the current branch to its existing tracking branch.
type PushCurrentBranchStep struct {
	CurrentBranch    domain.LocalBranchName
	InitialRemoteSHA domain.SHA
	NoPushHook       bool
	Undoable         bool
	EmptyStep
}

func (step *PushCurrentBranchStep) CreateUndoSteps(_ *git.BackendCommands) ([]Step, error) {
	if step.Undoable {
		return []Step{&ResetRemoteBranchToSHAStep{
			Branch:    step.CurrentBranch.RemoteName(),
			SHAToPush: step.InitialRemoteSHA,
		}}, nil
	}
	return []Step{&SkipCurrentBranchSteps{}}, nil
}

func (step *PushCurrentBranchStep) Run(run *git.ProdRunner, _ hosting.Connector) error {
	shouldPush, err := run.Backend.ShouldPushBranch(step.CurrentBranch, step.CurrentBranch.RemoteName())
	if err != nil {
		return err
	}
	if !shouldPush && !run.Config.DryRun {
		return nil
	}
	return run.Frontend.PushCurrentBranch(step.NoPushHook)
}
