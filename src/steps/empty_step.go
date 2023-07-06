package steps

import (
	"errors"

	"github.com/git-town/git-town/v9/src/git"
	"github.com/git-town/git-town/v9/src/hosting"
)

// EmptyStep does nothing.
// It is used for steps that have no undo or abort steps.
type EmptyStep struct{}

func (step *EmptyStep) CreateAbortStep() Step {
	return &EmptyStep{}
}

func (step *EmptyStep) CreateContinueStep() Step {
	return &EmptyStep{}
}

func (step *EmptyStep) CreateUndoStep(_ *git.BackendCommands) (Step, error) {
	return &EmptyStep{}, nil
}

func (step *EmptyStep) CreateAutomaticAbortError() error {
	return errors.New("")
}

func (step *EmptyStep) Run(_ *git.ProdRunner, connector hosting.Connector) error {
	return nil
}

func (step *EmptyStep) ShouldAutomaticallyAbortOnError() bool {
	return false
}
