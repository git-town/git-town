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
	Parent    string
	Force     bool
	branchSha string
}

func (step *DeleteLocalBranchStep) CreateUndoStep(repo *git.PublicRepo) (Step, error) {
	return &CreateBranchStep{Branch: step.Branch, StartingPoint: step.branchSha}, nil
}

func (step *DeleteLocalBranchStep) Run(repo *git.PublicRepo, connector hosting.Connector) error {
	var err error
	step.branchSha, err = repo.Internal.ShaForBranch(step.Branch)
	if err != nil {
		return err
	}
	hasUnmergedCommits, err := repo.Internal.BranchHasUnmergedCommits(step.Branch, step.Parent)
	if err != nil {
		return err
	}
	return repo.DeleteLocalBranch(step.Branch, step.Force || hasUnmergedCommits)
}
