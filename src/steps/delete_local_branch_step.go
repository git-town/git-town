package steps

import (
	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/hosting"
)

// DeleteLocalBranchStep deletes the branch with the given name,
// optionally in a safe or unsafe way.
type DeleteLocalBranchStep struct {
	NoOpStep
	BranchName string
	Force      bool
	branchSha  string
}

func (step *DeleteLocalBranchStep) CreateUndoStep(repo *git.ProdRepo) (Step, error) { //nolint:ireturn
	return &CreateBranchStep{BranchName: step.BranchName, StartingPoint: step.branchSha}, nil
}

func (step *DeleteLocalBranchStep) Run(repo *git.ProdRepo, driver hosting.Driver) error {
	var err error
	step.branchSha, err = repo.Silent.ShaForBranch(step.BranchName)
	if err != nil {
		return err
	}
	hasUnmergedCommits, err := repo.Silent.BranchHasUnmergedCommits(step.BranchName)
	if err != nil {
		return err
	}
	return repo.Logging.DeleteLocalBranch(step.BranchName, step.Force || hasUnmergedCommits)
}
