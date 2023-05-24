package steps

import (
	"github.com/git-town/git-town/v9/src/git"
	"github.com/git-town/git-town/v9/src/hosting"
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

func (step *DeleteLocalBranchStep) CreateUndoStep(backend *git.BackendCommands) (Step, error) {
	return &CreateBranchStep{Branch: step.Branch, StartingPoint: step.branchSha}, nil
}

func (step *DeleteLocalBranchStep) Run(run *git.ProdRunner, connector hosting.Connector) error {
	var err error
	step.branchSha, err = run.Backend.ShaForBranch(step.Branch)
	if err != nil {
		return err
	}
	hasUnmergedCommits, err := run.Backend.BranchHasUnmergedCommits(step.Branch, step.Parent)
	if err != nil {
		return err
	}
	return run.Frontend.DeleteLocalBranch(step.Branch, step.Force || hasUnmergedCommits)
}
