package steps

import (
	"errors"
	"fmt"

	"github.com/Originate/git-town/lib/git"
)

type EnsureHasShippableChangesStep struct {
	BranchName string
	NoUndoStep
}

func (step EnsureHasShippableChangesStep) CreateAbortStep() Step {
	return NoOpStep{}
}

func (step EnsureHasShippableChangesStep) CreateContinueStep() Step {
	return NoOpStep{}
}

func (step EnsureHasShippableChangesStep) GetAutomaticAbortErrorMessage() string {
	return fmt.Sprintf("The branch '%s' has no shippable changes.", step.BranchName)
}

func (step EnsureHasShippableChangesStep) Run() error {
	if !git.HasShippableChanges(step.BranchName) {
		return errors.New("no shippable changes")
	}
	return nil
}

func (step EnsureHasShippableChangesStep) ShouldAutomaticallyAbortOnError() bool {
	return true
}
