package steps

import (
	"github.com/git-town/git-town/v9/src/git"
	"github.com/git-town/git-town/v9/src/hosting"
)

// ContinueMergeStep finishes an ongoing merge conflict
// assuming all conflicts have been resolved by the user.
type ContinueMergeStep struct {
	EmptyStep
}

func (step *ContinueMergeStep) CreateContinueSteps() Step {
	return step
}

func (step *ContinueMergeStep) Run(run *git.ProdRunner, _ hosting.Connector) error {
	if run.Backend.HasMergeInProgress() {
		return run.Frontend.CommitNoEdit()
	}
	return nil
}
