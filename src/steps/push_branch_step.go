package steps

import (
	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/hosting"
)

// PushBranchStep pushes the branch with the given name to the origin remote.
// Optionally with force.
type PushBranchStep struct {
	NoOpStep
	BranchName     string
	Force          bool
	ForceWithLease bool
	NoPushVerify   bool
	Undoable       bool
}

func (step *PushBranchStep) CreateUndoStep(repo *git.ProdRepo) (Step, error) { //nolint:ireturn
	if step.Undoable {
		return &PushBranchAfterCurrentBranchSteps{}, nil
	}
	return &SkipCurrentBranchSteps{}, nil
}

func (step *PushBranchStep) Run(repo *git.ProdRepo, driver hosting.Driver) error {
	shouldPush, err := repo.Silent.ShouldPushBranch(step.BranchName)
	if err != nil {
		return err
	}
	if !shouldPush && !repo.DryRun.IsActive() {
		return nil
	}
	if step.ForceWithLease {
		return repo.Logging.PushBranchForceWithLease(step.BranchName, step.NoPushVerify)
	}
	if step.Force {
		return repo.Logging.PushBranchForce(step.BranchName, step.NoPushVerify)
	}
	currentBranch, err := repo.Silent.CurrentBranch()
	if err != nil {
		return err
	}
	if currentBranch == step.BranchName {
		return repo.Logging.PushBranch(step.NoPushVerify)
	}
	return repo.Logging.PushBranchToOrigin(step.BranchName, step.NoPushVerify)
}
