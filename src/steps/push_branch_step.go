package steps

import (
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/git"
	"github.com/git-town/git-town/v9/src/hosting"
)

// PushBranchStep pushes the branch with the given name to the origin remote.
// Optionally with force.
type PushBranchStep struct {
	Branch         domain.LocalBranchName
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
	shouldPush, err := run.Backend.ShouldPushBranch(step.Branch, step.Branch.AtRemote(domain.OriginRemote))
	if err != nil {
		return err
	}
	if !shouldPush && !run.Config.DryRun {
		return nil
	}
	return run.Frontend.PushBranch(git.PushArgs{
		Branch:         step.Branch,
		ForceWithLease: step.ForceWithLease,
		NoPushHook:     step.NoPushHook,
		Remote:         remote(step.Branch, step.Branch),
	})
}

// provides the name of the remote to push to.
func remote(currentBranch, stepBranch domain.LocalBranchName) domain.Remote {
	// TODO: how does this comparison of whether the branch in the step is the current branch make sense when deciding whether to push to origin or not?
	if currentBranch == stepBranch {
		return domain.NoRemote
	}
	return domain.OriginRemote
}
