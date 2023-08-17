package steps

import (
	"github.com/git-town/git-town/v9/src/config"
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/git"
	"github.com/git-town/git-town/v9/src/hosting"
)

// PushBranchStep pushes the branch with the given name to the origin remote.
// Optionally with force.
type PushBranchStep struct {
	EmptyStep
	Branch         domain.LocalBranchName
	TrackingBranch domain.RemoteBranchName
	ForceWithLease bool
	NoPushHook     bool
	Undoable       bool
}

func (step *PushBranchStep) CreateUndoSteps(_ *git.BackendCommands) ([]Step, error) {
	if step.Undoable {
		return []Step{&PushBranchAfterCurrentBranchSteps{}}, nil
	}
	return []Step{&SkipCurrentBranchSteps{}}, nil
}

func (step *PushBranchStep) Run(run *git.ProdRunner, _ hosting.Connector) error {
	shouldPush, err := run.Backend.ShouldPushBranch(step.Branch, step.TrackingBranch) // TODO: look this up in a git.Branches struct that needs to get injected here somehow
	if err != nil {
		return err
	}
	if !shouldPush && !run.Config.DryRun {
		return nil
	}
	currentBranch, err := run.Backend.CurrentBranch()
	if err != nil {
		return err
	}
	return run.Frontend.PushBranch(git.PushArgs{
		Branch:         step.Branch,
		ForceWithLease: step.ForceWithLease,
		NoPushHook:     step.NoPushHook,
		Remote:         remoteName(currentBranch, step.Branch),
	})
}

// provides the name of the remote to push to.
func remoteName(currentBranch, stepBranch domain.LocalBranchName) string {
	if currentBranch == stepBranch {
		return ""
	}
	return config.OriginRemote
}
