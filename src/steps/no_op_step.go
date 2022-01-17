//nolint:ireturn
package steps

import (
	"errors"

	"github.com/git-town/git-town/v7/src/drivers"
	"github.com/git-town/git-town/v7/src/git"
)

// NoOpStep does nothing.
// It is used for steps that have no undo or abort steps.
type NoOpStep struct{}

// CreateAbortStep returns the abort step for this step.
func (step *NoOpStep) CreateAbortStep() Step {
	return &NoOpStep{}
}

// CreateContinueStep returns the continue step for this step.
func (step *NoOpStep) CreateContinueStep() Step {
	return &NoOpStep{}
}

// CreateUndoStep returns the undo step for this step.
func (step *NoOpStep) CreateUndoStep(repo *git.ProdRepo) (Step, error) {
	return &NoOpStep{}, nil
}

// GetAutomaticAbortError returns the error message to display when this step
// cause the command to automatically abort.
func (step *NoOpStep) GetAutomaticAbortError() error {
	return errors.New("")
}

// Run executes this step.
func (step *NoOpStep) Run(repo *git.ProdRepo, driver drivers.CodeHostingDriver) error {
	return nil
}

// ShouldAutomaticallyAbortOnError returns whether this step should cause the command to
// automatically abort if it errors.
func (step *NoOpStep) ShouldAutomaticallyAbortOnError() bool {
	return false
}
