package steps

import (
	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/hosting"
)

// SetParentBranchStep registers the branch with the given name as a parent
// of the branch with the other given name.
type SetParentBranchStep struct {
	NoOpStep
	Branch         string
	ParentBranch   string
	previousParent string
}

func (step *SetParentBranchStep) CreateUndoStep(repo *git.ProdRepo) (Step, error) {
	if step.previousParent == "" {
		return &DeleteParentBranchStep{Branch: step.Branch}, nil
	}
	return &SetParentBranchStep{Branch: step.Branch, ParentBranch: step.previousParent}, nil
}

func (step *SetParentBranchStep) Run(repo *git.ProdRepo, connector hosting.Connector) error {
	step.previousParent = repo.Config.ParentBranch(step.Branch)
	return repo.Config.SetParentBranch(step.Branch, step.ParentBranch)
}
