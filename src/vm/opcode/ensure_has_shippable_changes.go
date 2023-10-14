package opcode

import (
	"fmt"

	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/messages"
)

// EnsureHasShippableChanges asserts that the branch has unique changes not on the main branch.
type EnsureHasShippableChanges struct {
	Branch domain.LocalBranchName
	Parent domain.LocalBranchName
	BaseOpcode
}

func (step *EnsureHasShippableChanges) CreateAutomaticAbortError() error {
	return fmt.Errorf(messages.ShipBranchNothingToDo, step.Branch)
}

func (step *EnsureHasShippableChanges) Run(args RunArgs) error {
	hasShippableChanges, err := args.Runner.Backend.HasShippableChanges(step.Branch, step.Parent)
	if err != nil {
		return err
	}
	if !hasShippableChanges {
		return fmt.Errorf(messages.ShipBranchNothingToDo, step.Branch)
	}
	return nil
}

func (step *EnsureHasShippableChanges) ShouldAutomaticallyAbortOnError() bool {
	return true
}
