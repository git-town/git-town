package steps

import (
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/git"
	"github.com/git-town/git-town/v9/src/hosting"
)

// PushBranchStep pushes the branch with the given name to the origin remote.
// Optionally with force.
type ForcePushBranchStep struct {
	Branch     domain.LocalBranchName
	NoPushHook bool
	EmptyStep
}

func (step *ForcePushBranchStep) CreateUndoSteps(_ *git.BackendCommands) ([]Step, error) {
	return []Step{&SkipCurrentBranchSteps{}}, nil
}

func (step *ForcePushBranchStep) Run(run *git.ProdRunner, _ hosting.Connector) error {
	shouldPush, err := run.Backend.ShouldPushBranch(step.Branch, step.Branch.RemoteName()) // TODO: look this up in a git.Branches struct that needs to get injected here somehow
	if err != nil {
		return err
	}
	if !shouldPush && !run.Config.DryRun {
		return nil
	}
	return run.Frontend.ForcePushBranch(step.NoPushHook)
}
