package opcode

import (
	"fmt"

	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/messages"
	"github.com/git-town/git-town/v9/src/vm/shared"
)

// EnsureHasShippableChanges asserts that the branch has unique changes not on the main branch.
type EnsureHasShippableChanges struct {
	Branch domain.LocalBranchName
	Parent domain.LocalBranchName
	undeclaredOpcodeMethods
}

func (op *EnsureHasShippableChanges) CreateAutomaticAbortError() error {
	return fmt.Errorf(messages.ShipBranchNothingToDo, op.Branch)
}

func (op *EnsureHasShippableChanges) Run(args shared.RunArgs) error {
	hasShippableChanges, err := args.Runner.Backend.HasShippableChanges(op.Branch, op.Parent)
	if err != nil {
		return err
	}
	if !hasShippableChanges {
		return fmt.Errorf(messages.ShipBranchNothingToDo, op.Branch)
	}
	return nil
}

func (op *EnsureHasShippableChanges) ShouldAutomaticallyAbortOnError() bool {
	return true
}
