package steps

import (
	"errors"
)

// EmptyStep does nothing.
// It is used for steps that have no undo or abort steps.
type EmptyStep struct{}

func (step *EmptyStep) CreateAbortSteps() []Step {
	return []Step{}
}

func (step *EmptyStep) CreateContinueSteps() []Step {
	return []Step{}
}

func (step *EmptyStep) CreateAutomaticAbortError() error {
	return errors.New("")
}

func (step *EmptyStep) Run(_ RunArgs) error {
	return nil
}

func (step *EmptyStep) ShouldAutomaticallyAbortOnError() bool {
	return false
}
