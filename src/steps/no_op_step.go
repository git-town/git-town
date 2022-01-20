package steps

import (
	"errors"

	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/hosting"
)

// NoOpStep does nothing.
// It is used for steps that have no undo or abort steps.
type NoOpStep struct{}

func (step *NoOpStep) CreateAbortStep() Step { //nolint:ireturn
	return &NoOpStep{}
}

func (step *NoOpStep) CreateContinueStep() Step { //nolint:ireturn
	return &NoOpStep{}
}

func (step *NoOpStep) CreateUndoStep(repo *git.ProdRepo) (Step, error) { //nolint:ireturn
	return &NoOpStep{}, nil
}

func (step *NoOpStep) CreateAutomaticAbortError() error {
	return errors.New("")
}

func (step *NoOpStep) Run(repo *git.ProdRepo, driver hosting.Driver) error {
	return nil
}

func (step *NoOpStep) ShouldAutomaticallyAbortOnError() bool {
	return false
}
