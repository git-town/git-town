package steps

import (
	"github.com/git-town/git-town/v7/src/config"
	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/hosting"
)

// PushBranchStep pushes the branch with the given name to the origin remote.
// Optionally with force.
type PushBranchStep struct {
	EmptyStep
	Branch         string
	Force          bool
	ForceWithLease bool
	NoPushHook     bool
	Undoable       bool
}

func (step *PushBranchStep) CreateUndoStep(repo *git.ProdRepo) (Step, error) {
	if step.Undoable {
		return &PushBranchAfterCurrentBranchSteps{}, nil
	}
	return &SkipCurrentBranchSteps{}, nil
}

func (step *PushBranchStep) Run(repo *git.ProdRepo, connector hosting.Connector) error {
	shouldPush, err := repo.Silent.ShouldPushBranch(step.Branch)
	if err != nil {
		return err
	}
	if !shouldPush && !repo.DryRun.IsActive() {
		return nil
	}
	currentBranch, err := repo.Silent.CurrentBranch()
	if err != nil {
		return err
	}
	return repo.Logging.PushBranch(git.PushArgs{
		Branch:         step.Branch,
		ForceWithLease: step.ForceWithLease,
		NoPushHook:     step.NoPushHook,
		Force:          step.Force,
		Remote:         remoteName(currentBranch, step.Branch),
	})
}

// provides the name of the remote to push to.
func remoteName(currentBranch, stepBranch string) string {
	if currentBranch == stepBranch {
		return ""
	}
	return config.OriginRemote
}
