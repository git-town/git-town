package steps

import (
	"github.com/git-town/git-town/v9/src/git"
	"github.com/git-town/git-town/v9/src/hosting"
)

// ContinueRebaseStep finishes an ongoing rebase operation
// assuming all conflicts have been resolved by the user.
type ContinueRebaseStep struct {
	EmptyStep
}

func (step *ContinueRebaseStep) CreateAbortSteps() []Step {
	return []Step{&AbortRebaseStep{}}
}

func (step *ContinueRebaseStep) CreateContinueSteps() []Step {
	return []Step{step}
}

func (step *ContinueRebaseStep) Run(run *git.ProdRunner, _ hosting.Connector) error {
	hasRebaseInProgress, err := run.Backend.HasRebaseInProgress()
	if err != nil {
		return err
	}
	if hasRebaseInProgress {
		return run.Frontend.ContinueRebase()
	}
	return nil
}
