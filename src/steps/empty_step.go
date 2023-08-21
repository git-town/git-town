package steps

import (
	"errors"

	"github.com/git-town/git-town/v9/src/git"
	"github.com/git-town/git-town/v9/src/hosting"
)

// EmptyStep does nothing.
// It is used for steps that have no undo or abort steps.
type EmptyStep struct{}

func (step *EmptyStep) CreateAbortSteps() Step {
	return &EmptyStep{}
}

func (step *EmptyStep) CreateContinueSteps() Step {
	return &EmptyStep{}
}

func (step *EmptyStep) CreateUndoSteps(_ *git.BackendCommands) ([]Step, error) {
	return []Step{}, nil
}

func (step *EmptyStep) CreateAutomaticAbortError() error {
	return errors.New("")
}

func (step *EmptyStep) Run(_ *git.ProdRunner, _ hosting.Connector) error {
	return nil
}

func (step *EmptyStep) ShouldAutomaticallyAbortOnError() bool {
	return false
}
