package steps

import (
	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/hosting"
)

// DeleteLocalBranchStep deletes the branch with the given name,
// optionally in a safe or unsafe way.
type DeleteLocalBranchStep struct {
	EmptyStep
	Branch    string
	Force     bool
	branchSha string
}

func (step *DeleteLocalBranchStep) CreateUndoStep(repo *git.ProdRepo) (Step, error) {
	return &CreateBranchStep{Branch: step.Branch, StartingPoint: step.branchSha}, nil
}

func (step *DeleteLocalBranchStep) Run(repo *git.ProdRepo, connector hosting.Connector) error {
	var err error
	step.branchSha, err = repo.Silent.ShaForBranch(step.Branch)
	if err != nil {
		return err
	}
	hasUnmergedCommits, err := repo.Silent.BranchHasUnmergedCommits(step.Branch)
	if err != nil {
		return err
	}
	return repo.Logging.DeleteLocalBranch(step.Branch, step.Force || hasUnmergedCommits)
}
