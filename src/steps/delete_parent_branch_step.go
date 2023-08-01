package steps

import (
	"github.com/git-town/git-town/v9/src/git"
	"github.com/git-town/git-town/v9/src/hosting"
)

// DeleteParentBranchStep removes the parent branch entry in the Git Town configuration.
type DeleteParentBranchStep struct {
	EmptyStep
	Branch string
	Parent string
}

func (step *DeleteParentBranchStep) CreateUndoSteps(_ *git.BackendCommands) ([]Step, error) {
	if step.Parent == "" {
		return []Step{}, nil
	}
	return []Step{&SetParentStep{Branch: step.Branch, ParentBranch: step.Parent}}, nil
}

func (step *DeleteParentBranchStep) Run(run *git.ProdRunner, _ hosting.Connector) error {
	return run.Config.RemoveParent(step.Branch)
}
