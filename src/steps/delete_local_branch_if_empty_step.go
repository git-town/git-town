package steps

import (
	"fmt"

	"github.com/git-town/git-town/v9/src/git"
	"github.com/git-town/git-town/v9/src/hosting"
	"github.com/git-town/git-town/v9/src/messages"
)

// DeleteLocalBranchStep deletes the branch with the given name,
// optionally in a safe or unsafe way.
type DeleteLocalBranchIfEmptyStep struct {
	EmptyStep
	Branch        string
	Parent        string
	Force         bool
	branchSha     string
	branchDeleted bool
}

func (step *DeleteLocalBranchIfEmptyStep) CreateUndoStep(_ *git.BackendCommands) ([]Step, error) {
	if !step.branchDeleted {
		return []Step{}, nil
	}
	return []Step{
		&CreateBranchStep{
			Branch:        step.Branch,
			StartingPoint: step.branchSha,
		},
		&SetParentStep{
			Branch:       step.Branch,
			ParentBranch: step.Parent,
		},
	}, nil
}

func (step *DeleteLocalBranchIfEmptyStep) Run(run *git.ProdRunner, _ hosting.Connector) error {
	var err error
	step.branchSha, err = run.Backend.ShaForBranch(step.Branch)
	if err != nil {
		return err
	}
	// ensure branch is empty
	branchHasUnmergedChanges, err := run.Backend.BranchHasUnmergedCommits(step.Branch, step.Parent)
	if err != nil {
		return err
	}
	if branchHasUnmergedChanges {
		run.Stats.RegisterMessage(fmt.Sprintf(messages.BranchDeletedHasUnmergedChanges, step.Branch))
		return nil
	}
	// delete the local branch
	err = run.Frontend.DeleteLocalBranch(step.Branch, step.Force)
	if err != nil {
		return err
	}
	step.branchDeleted = true
	// delete the configuration settings for this branch
	err = run.Backend.Config.RemoveParent(step.Branch)
	if err != nil {
		return err
	}
	// updating the proposals of the child branches is not necessary here since the remote branch was already deleted
	// so there shouldn't be any open proposals
	return nil
}
