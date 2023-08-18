package steps

import (
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/git"
	"github.com/git-town/git-town/v9/src/hosting"
)

// DeleteLocalBranchStep deletes the branch with the given name,
// optionally in a safe or unsafe way.
type DeleteLocalBranchStep struct {
	EmptyStep
	Branch    domain.LocalBranchName
	Parent    domain.Location
	Force     bool
	branchSha domain.SHA
}

func (step *DeleteLocalBranchStep) CreateUndoSteps(_ *git.BackendCommands) ([]Step, error) {
	return []Step{&CreateBranchStep{Branch: step.Branch, StartingPoint: step.branchSha.Location()}}, nil
}

func (step *DeleteLocalBranchStep) Run(run *git.ProdRunner, _ hosting.Connector) error {
	var err error
	step.branchSha, err = run.Backend.ShaForBranch(step.Branch.BranchName())
	if err != nil {
		return err
	}
	hasUnmergedCommits, err := run.Backend.BranchHasUnmergedCommits(step.Branch, step.Parent)
	if err != nil {
		return err
	}
	return run.Frontend.DeleteLocalBranch(step.Branch, step.Force || hasUnmergedCommits)
}
