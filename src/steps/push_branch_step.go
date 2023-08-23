package steps

import (
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/git"
	"github.com/git-town/git-town/v9/src/hosting"
)

// PushBranchStep pushes the branch with the given name to the origin remote.
// Optionally with force.
type PushBranchStep struct {
	Branch domain.LocalBranchName
	// TrackingBranch domain.RemoteBranchName // TODO: populate this with the actual tracking branch name
	ForceWithLease bool
	NoPushHook     bool
	Undoable       bool
	EmptyStep
}

func (step *PushBranchStep) CreateUndoSteps(_ *git.BackendCommands) ([]Step, error) {
	if step.Undoable {
		return []Step{&PushBranchAfterCurrentBranchSteps{}}, nil
	}
	return []Step{&SkipCurrentBranchSteps{}}, nil
}

func (step *PushBranchStep) Run(run *git.ProdRunner, _ hosting.Connector) error {
	shouldPush, err := run.Backend.ShouldPushBranch(step.Branch, step.Branch.RemoteName()) // TODO: look this up in a git.Branches struct that needs to get injected here somehow
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
		Remote:         remote(currentBranch, step.Branch),
	})
}

// provides the name of the remote to push to.
func remote(currentBranch, stepBranch domain.LocalBranchName) domain.Remote {
	if currentBranch == stepBranch {
		return domain.NoRemote
	}
	return domain.OriginRemote
}
