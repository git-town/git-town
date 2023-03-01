package steps

import (
	"errors"
	"fmt"

	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/hosting"
)

// EnsureHasShippableChangesStep asserts that the branch has unique changes not on the main branch.
type EnsureHasShippableChangesStep struct {
	EmptyStep
	Branch string
}

func (step *EnsureHasShippableChangesStep) CreateAutomaticAbortError() error {
	return fmt.Errorf("the branch %q has no shippable changes", step.Branch)
}

func (step *EnsureHasShippableChangesStep) Run(repo *git.ProdRepo, connector hosting.Connector) error {
	hasShippableChanges, err := repo.Silent.HasShippableChanges(step.Branch)
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
