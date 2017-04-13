package steps

import (
	"errors"
	"fmt"

	"github.com/Originate/git-town/lib/git"
)

// EnsureHasShippableChangesStep asserts that the branch has unique changes not on the main branch
type EnsureHasShippableChangesStep struct {
	NoAbortStep
	NoContinueStep
	NoUndoStep
	BranchName string
}

// GetAutomaticAbortErrorMessage returns the error message to display when this step
// cause the command to automatically abort.
func (step EnsureHasShippableChangesStep) GetAutomaticAbortErrorMessage() string {
	return fmt.Sprintf("The branch '%s' has no shippable changes.", step.BranchName)
}

// Run executes this step.
func (step EnsureHasShippableChangesStep) Run() error {
	if !git.HasShippableChanges(step.BranchName) {
		return errors.New("no shippable changes")
	}
	return nil
}

// ShouldAutomaticallyAbortOnError returns whether this step should cause the command to
// automatically abort if it errors.
func (step EnsureHasShippableChangesStep) ShouldAutomaticallyAbortOnError() bool {
	return true
}
