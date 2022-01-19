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

func (step *EnsureHasShippableChangesStep) CreateAutomaticAbortError() error {
	return fmt.Errorf("the branch %q has no shippable changes", step.BranchName)
}

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

func (step *EnsureHasShippableChangesStep) ShouldAutomaticallyAbortOnError() bool {
	return true
}
