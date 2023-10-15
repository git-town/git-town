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

func (self *EnsureHasShippableChanges) CreateAutomaticAbortError() error {
	return fmt.Errorf(messages.ShipBranchNothingToDo, self.Branch)
}

func (self *EnsureHasShippableChanges) Run(args shared.RunArgs) error {
	hasShippableChanges, err := args.Runner.Backend.HasShippableChanges(self.Branch, self.Parent)
	if err != nil {
		return err
	}
	if !hasShippableChanges {
		return fmt.Errorf(messages.ShipBranchNothingToDo, self.Branch)
	}
	return nil
}

func (self *EnsureHasShippableChanges) ShouldAutomaticallyAbortOnError() bool {
	return true
}
