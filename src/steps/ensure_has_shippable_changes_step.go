package steps

import (
	"errors"
	"fmt"

	"github.com/git-town/git-town/v7/src/drivers"
	"github.com/git-town/git-town/v7/src/git"
)

// EnsureHasShippableChangesStep asserts that the branch has unique changes not on the main branch.
type EnsureHasShippableChangesStep struct {
	NoOpStep
	BranchName string
}

// CreateAutomaticAbortError returns the error message to display when this step
// cause the command to automatically abort.
func (step *EnsureHasShippableChangesStep) CreateAutomaticAbortError() error {
	return fmt.Errorf("the branch %q has no shippable changes", step.BranchName)
}

// Run executes this step.
func (step *EnsureHasShippableChangesStep) Run(repo *git.ProdRepo, driver drivers.CodeHostingDriver) error {
	hasShippableChanges, err := repo.Silent.HasShippableChanges(step.BranchName)
	if err != nil {
		return err
	}
	if !hasShippableChanges {
		return errors.New("no shippable changes")
	}
	return nil
}

// ShouldAutomaticallyAbortOnError returns whether this step should cause the command to
// automatically abort if it errors.
func (step *EnsureHasShippableChangesStep) ShouldAutomaticallyAbortOnError() bool {
	return true
}
