package step

import (
	"github.com/git-town/git-town/v9/src/domain"
)

// PushCurrentBranch pushes the current branch to its existing tracking branch.
type PushCurrentBranch struct {
	CurrentBranch domain.LocalBranchName
	NoPushHook    bool
	Empty
}

func (step *PushCurrentBranch) Run(args RunArgs) error {
	shouldPush, err := args.Runner.Backend.ShouldPushBranch(step.CurrentBranch, step.CurrentBranch.TrackingBranch())
	if err != nil {
		return err
	}
	if !shouldPush && !args.Runner.Config.DryRun {
		return nil
	}
	return args.Runner.Frontend.PushCurrentBranch(step.NoPushHook)
}
