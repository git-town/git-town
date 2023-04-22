package steps

import (
	"github.com/git-town/git-town/v8/src/git"
	"github.com/git-town/git-town/v8/src/hosting"
)

// CreateBranchStep creates a new branch
// but leaves the current branch unchanged.
type CreateBranchStep struct {
	EmptyStep
	Branch        string
	StartingPoint string
}

func (step *CreateBranchStep) CreateUndoStep(backend *git.BackendCommands) (Step, error) {
	return &DeleteLocalBranchStep{Branch: step.Branch, Parent: step.StartingPoint, Force: true}, nil
}

func (step *CreateBranchStep) Run(run *git.ProdRunner, connector hosting.Connector) error {
	return run.Frontend.CreateBranch(step.Branch, step.StartingPoint)
}
