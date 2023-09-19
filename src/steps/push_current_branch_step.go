package steps

import (
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/git"
)

// PushCurrentBranchStep pushes the current branch to its tracking branch.
// The tracking branch must exist.
type PushCurrentBranchStep struct {
	CurrentBranch    domain.LocalBranchName
	initialRemoteSHA domain.SHA
	shaAfterPush     domain.SHA
	NoPushHook       bool
	Undoable         bool
	EmptyStep
}

func (step *PushCurrentBranchStep) CreateUndoSteps(_ *git.BackendCommands) ([]Step, error) {
	if step.Undoable {
		return []Step{&ResetRemoteBranchToSHAStep{
			Branch:           step.CurrentBranch.RemoteBranch(),
			SHAToPush:        step.initialRemoteSHA,
			SHAThatMustExist: step.shaAfterPush,
		}}, nil
	}
	return []Step{&SkipCurrentBranchSteps{}}, nil
}

func (step *PushCurrentBranchStep) Run(args RunArgs) error {
	trackingBranch := step.CurrentBranch.RemoteBranch()
	shouldPush, err := run.Backend.ShouldPushBranch(step.CurrentBranch, trackingBranch)
	if err != nil {
		return err
	}
	if !shouldPush && !args.Runner.Config.DryRun {
		return nil
	}
	step.initialRemoteSHA, err = run.Backend.SHAForBranch(trackingBranch.BranchName())
	if err != nil {
		return err
	}
	step.shaAfterPush, err = run.Backend.CurrentSHA()
	if err != nil {
		return err
	}
	return args.Runner.Frontend.PushCurrentBranch(step.NoPushHook)
}
