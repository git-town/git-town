package steps

import (
	"github.com/git-town/git-town/src/git"
)

// DeleteLocalBranchStep deletes the branch with the given name,
// optionally in a safe or unsafe way.
type DeleteLocalBranchStep struct {
	NoOpStep
	BranchName string
	Force      bool

	branchSha string
}

// CreateUndoStep returns the undo step for this step.
func (step *DeleteLocalBranchStep) CreateUndoStep() Step {
	return &CreateBranchStep{BranchName: step.BranchName, StartingPoint: step.branchSha}
}

// Run executes this step.
func (step *DeleteLocalBranchStep) Run(repo *git.ProdRepo) (err error) {
	step.branchSha, err = repo.Silent.BranchSha(step.BranchName)
	if err != nil {
		return err
	}
	hasUnmergedCommits, err := repo.Silent.BranchHasUnmergedCommits(step.BranchName)
	if err != nil {
		return err
	}
	return repo.Logging.DeleteLocalBranch(step.BranchName, step.Force || hasUnmergedCommits)
}
