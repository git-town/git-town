package steps

import (
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/git"
)

// PushCurrentBranchStep pushes the current branch to its existing tracking branch.
type PushCurrentBranchStep struct {
	CurrentBranch domain.LocalBranchName
	NoPushHook    bool
	Undoable      bool
	EmptyStep
}

func (step *PushCurrentBranchStep) CreateUndoSteps(_ *git.BackendCommands) ([]Step, error) {
	if step.Undoable {
		return []Step{&PushBranchAfterCurrentBranchSteps{}}, nil
	}
	return []Step{&SkipCurrentBranchSteps{}}, nil
}

func (step *PushCurrentBranchStep) Run(args RunArgs) error {
	shouldPush, err := args.Run.Backend.ShouldPushBranch(step.CurrentBranch, step.CurrentBranch.RemoteBranch())
	if err != nil {
		return err
	}
	if !shouldPush && !args.Run.Config.DryRun {
		return nil
	}
	return args.Run.Frontend.PushCurrentBranch(step.NoPushHook)
}
