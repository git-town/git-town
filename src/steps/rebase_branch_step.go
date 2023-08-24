package steps

import (
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/git"
	"github.com/git-town/git-town/v9/src/hosting"
)

// RebaseBranchStep rebases the current branch
// against the branch with the given name.
type RebaseBranchStep struct {
	Branch      domain.BranchName
	previousSHA domain.SHA
	EmptyStep
}

func (step *RebaseBranchStep) CreateAbortSteps() []Step {
	return []Step{&AbortRebaseStep{}}
}

func (step *RebaseBranchStep) CreateContinueSteps() []Step {
	return []Step{&ContinueRebaseStep{}}
}

func (step *RebaseBranchStep) CreateUndoSteps(_ *git.BackendCommands) ([]Step, error) {
	return []Step{&ResetToSHAStep{Hard: true, SHA: step.previousSHA}}, nil
}

func (step *RebaseBranchStep) Run(run *git.ProdRunner, _ hosting.Connector) error {
	var err error
	step.previousSHA, err = run.Backend.CurrentSHA()
	if err != nil {
		return err
	}
	return run.Frontend.Rebase(step.Branch)
}
