package steps

import (
	"github.com/git-town/git-town/v8/src/git"
	"github.com/git-town/git-town/v8/src/hosting"
)

// DeleteParentBranchStep removes the parent branch entry in the Git Town configuration.
type DeleteParentBranchStep struct {
	EmptyStep
	Branch         string
	previousParent string // TODO: make public and populate as part of creating this step
}

func (step *DeleteParentBranchStep) CreateUndoStep(backend *git.BackendCommands) (Step, error) {
	if step.previousParent == "" {
		return &EmptyStep{}, nil
	}
	return &SetParentStep{Branch: step.Branch, ParentBranch: step.previousParent}, nil
}

func (step *DeleteParentBranchStep) Run(run *git.ProdRunner, connector hosting.Connector) error {
	step.previousParent = run.Config.ParentBranch(step.Branch)
	return run.Config.RemoveParent(step.Branch)
}
