package steps

import (
	"errors"

	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/hosting"
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

func (step *EmptyStep) CreateUndoStep(repo *git.ProdRepo) (Step, error) {
	return &EmptyStep{}, nil
}

func (step *EmptyStep) CreateAutomaticAbortError() error {
	return errors.New("")
}

func (step *EmptyStep) Run(repo *git.ProdRepo, connector hosting.Connector) error {
	return nil
}

func (step *EmptyStep) ShouldAutomaticallyAbortOnError() bool {
	return false
}
