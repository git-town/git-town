package steps

import (
	"fmt"

	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/git"
	"github.com/git-town/git-town/v9/src/hosting"
	"github.com/git-town/git-town/v9/src/messages"
)

// EnsureHasShippableChangesStep asserts that the branch has unique changes not on the main branch.
type EnsureHasShippableChangesStep struct {
	EmptyStep
	Branch domain.LocalBranchName
	Parent domain.LocalBranchName
}

func (step *EnsureHasShippableChangesStep) CreateAutomaticAbortError() error {
	return fmt.Errorf(messages.ShipBranchNothingToDo, step.Branch)
}

func (step *EnsureHasShippableChangesStep) Run(run *git.ProdRunner, _ hosting.Connector) error {
	hasShippableChanges, err := run.Backend.HasShippableChanges(step.Branch, step.Parent)
	if err != nil {
		return err
	}
	if !hasShippableChanges {
		return fmt.Errorf(messages.ShipBranchNothingToDo, step.Branch)
	}
	return nil
}

func (step *EnsureHasShippableChangesStep) ShouldAutomaticallyAbortOnError() bool {
	return true
}
