package steps

import (
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/git"
)

// ForcePushBranchStep force-pushes the branch with the given name to the origin remote.
type ForcePushBranchStep struct {
	Branch     domain.LocalBranchName
	NoPushHook bool
	EmptyStep
}

func (step *ForcePushBranchStep) CreateUndoSteps(_ *git.BackendCommands) ([]Step, error) {
	return []Step{&SkipCurrentBranchSteps{}}, nil
}

func (step *ForcePushBranchStep) Run(args RunArgs) error {
	shouldPush, err := args.Run.Backend.ShouldPushBranch(step.Branch, step.Branch.RemoteBranch())
	if err != nil {
		return err
	}
	if !shouldPush && !args.Run.Config.DryRun {
		return nil
	}
	return args.Run.Frontend.ForcePushBranch(step.NoPushHook)
}
