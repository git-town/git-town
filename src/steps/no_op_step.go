package steps

import (
	"errors"

	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/hosting"
)

// NoOpStep does nothing.
// It is used for steps that have no undo or abort steps.
type NoOpStep struct{}

func (step *NoOpStep) CreateAbortStep() Step {
	return &NoOpStep{}
}

func (step *NoOpStep) CreateContinueStep() Step {
	return &NoOpStep{}
}

func (step *NoOpStep) CreateUndoStep(repo *git.ProdRepo) (Step, error) {
	return &NoOpStep{}, nil
}

func (step *NoOpStep) CreateAutomaticAbortError() error {
	return errors.New("")
}

func (step *NoOpStep) Run(repo *git.ProdRepo, connector hosting.Connector) error {
	return nil
}

func (step *NoOpStep) ShouldAutomaticallyAbortOnError() bool {
	return false
}
