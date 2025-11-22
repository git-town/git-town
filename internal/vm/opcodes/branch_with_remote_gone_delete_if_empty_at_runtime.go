package opcodes

import (
	"fmt"

	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/messages"
	"github.com/git-town/git-town/v22/internal/vm/shared"
)

// BranchWithRemoteGoneDeleteIfEmptyAtRuntime deletes the given branch if it has no content at runtime.
type BranchWithRemoteGoneDeleteIfEmptyAtRuntime struct {
	Branch gitdomain.LocalBranchName
}

func (self *BranchWithRemoteGoneDeleteIfEmptyAtRuntime) Run(args shared.RunArgs) error {
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
			&CheckoutAncestorOrOtherIfNeeded{
				Branch: self.Branch,
			},
			&BranchLocalDeleteContent{
				BranchToDelete:     self.Branch,
				BranchToRebaseOnto: args.Config.Value.ValidatedConfigData.MainBranch,
			},
			&LineageBranchRemove{
				Branch: self.Branch,
			},
		)
	}
	return nil
}
