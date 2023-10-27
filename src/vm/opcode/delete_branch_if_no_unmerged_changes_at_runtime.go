package opcode

import (
	"fmt"

	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/messages"
	"github.com/git-town/git-town/v9/src/vm/shared"
)

// DeleteBranchIfNoUnmergedChangesAtRuntime allows running different opcodes based on a condition evaluated at runtime.
type DeleteBranchIfNoUnmergedChangesAtRuntime struct {
	Branch domain.LocalBranchName
	undeclaredOpcodeMethods
}

func (self *DeleteBranchIfNoUnmergedChangesAtRuntime) Run(args shared.RunArgs) error {
	parent := args.Lineage.Parent(self.Branch)
	hasUnmergedChanges, err := args.Runner.Backend.BranchHasUnmergedChanges(self.Branch, parent)
	if err != nil {
		return err
	}
	if hasUnmergedChanges {
		args.PrependOpcodes(&QueueMessage{
			Message: fmt.Sprintf(messages.BranchDeletedHasUnmergedChanges, self.Branch),
		})
	} else {
		args.PrependOpcodes(
			&CheckoutParent{CurrentBranch: self.Branch},
			&DeleteLocalBranch{
				Branch: self.Branch,
				Force:  false,
			},
			&RemoveBranchFromLineage{
				Branch: self.Branch,
			},
			&QueueMessage{
				Message: fmt.Sprintf(messages.BranchDeleted, self.Branch),
			})
	}
	return nil
}
