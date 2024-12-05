package opcodes

import (
	"fmt"

	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/messages"
	"github.com/git-town/git-town/v16/internal/vm/shared"
)

// BranchDeleteIfEmptyAtRuntime allows running different opcodes based on a condition evaluated at runtime.
type BranchDeleteIfEmptyAtRuntime struct {
	Branch                  gitdomain.LocalBranchName
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *BranchDeleteIfEmptyAtRuntime) Run(args shared.RunArgs) error {
	parent, hasParent := args.Config.Value.NormalConfig.Lineage.Parent(self.Branch).Get()
	if !hasParent {
		return nil
	}
	hasUnmergedChanges, err := args.Git.BranchHasUnmergedChanges(args.Backend, self.Branch, parent)
	if err != nil {
		return err
	}
	if hasUnmergedChanges {
		args.PrependOpcodes(&MessageQueue{
			Message: fmt.Sprintf(messages.BranchDeletedHasUnmergedChanges, self.Branch),
		})
	} else {
		args.PrependOpcodes(
			&CheckoutParentOrMain{Branch: self.Branch},
			&BranchLocalDelete{Branch: self.Branch},
			&LineageBranchRemove{
				Branch: self.Branch,
			},
		)
	}
	return nil
}
