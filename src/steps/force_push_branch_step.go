package steps

import (
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/git"
)

// ForcePushBranchStep force-pushes the branch with the given name to the origin remote.
// TODO: rename to ForcePushCurrentBranchStep and determine the current branch
type ForcePushBranchStep struct {
	Branch     domain.LocalBranchName
	NoPushHook bool
	EmptyStep
}

func (step *ForcePushBranchStep) CreateUndoSteps(_ *git.BackendCommands) ([]Step, error) {
	return []Step{&SkipCurrentBranchSteps{}}, nil
}

func (step *ForcePushBranchStep) Run(args RunArgs) error {
	shouldPush, err := args.Runner.Backend.ShouldPushBranch(step.Branch, step.Branch.RemoteBranch())
	if err != nil {
		return err
	}
	if !shouldPush && !args.Runner.Config.DryRun {
		return nil
	}
	return args.Runner.Frontend.ForcePushBranch(step.NoPushHook)
}
