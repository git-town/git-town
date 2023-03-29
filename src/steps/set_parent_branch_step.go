package steps

import (
	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/hosting"
)

// SetParentStep registers the branch with the given name as a parent
// of the branch with the other given name.
type SetParentStep struct {
	EmptyStep
	Branch         string
	ParentBranch   string
	previousParent string
}

func (step *SetParentStep) CreateUndoStep(backend *git.BackendCommands) (Step, error) {
	if step.previousParent == "" {
		return &DeleteParentBranchStep{Branch: step.Branch}, nil
	}
	return &SetParentStep{Branch: step.Branch, ParentBranch: step.previousParent}, nil
}

func (step *SetParentStep) Run(run *git.ProdRunner, connector hosting.Connector) error {
	step.previousParent = run.Config.ParentBranch(step.Branch)
	return run.Config.SetParent(step.Branch, step.ParentBranch)
}
