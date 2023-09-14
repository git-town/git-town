package steps

import (
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/git"
)

// CreateBranchStep creates a new branch but leaves the current branch unchanged.
type CreateBranchStep struct {
	Branch        domain.LocalBranchName
	StartingPoint domain.Location
	EmptyStep
}

func (step *CreateBranchStep) CreateUndoSteps(_ *git.BackendCommands) ([]Step, error) {
	return []Step{&DeleteLocalBranchStep{Branch: step.Branch, Parent: step.StartingPoint, Force: true}}, nil
}

func (step *CreateBranchStep) Run(args RunArgs) error {
	return args.Run.Frontend.CreateBranch(step.Branch, step.StartingPoint)
}
