package steps

import (
	"github.com/git-town/git-town/v9/src/domain"
)

// PushCurrentBranchStep pushes the current branch to its existing tracking branch.
type PushCurrentBranchStep struct {
	CurrentBranch domain.LocalBranchName
	NoPushHook    bool
	EmptyStep
}

func (step *PushCurrentBranchStep) Run(args RunArgs) error {
	shouldPush, err := args.Runner.Backend.ShouldPushBranch(step.CurrentBranch, step.CurrentBranch.RemoteBranch())
	if err != nil {
		return err
	}
	if !shouldPush && !args.Runner.Config.DryRun {
		return nil
	}
	return args.Runner.Frontend.PushCurrentBranch(step.NoPushHook)
}
