package opcodes

import (
	"fmt"

	"github.com/git-town/git-town/v22/internal/forge/forgedomain"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/vm/shared"
)

// ProposalUpdateTargetToGrandParent updates the target of the proposal with the given number to the parent of the given branch.
type ProposalUpdateTargetToGrandParent struct {
	Branch    gitdomain.LocalBranchName
	OldTarget gitdomain.LocalBranchName
	Proposal  forgedomain.Proposal
}

func (self *ProposalUpdateTargetToGrandParent) Run(args shared.RunArgs) error {
	parent, hasParent := args.Config.Value.NormalConfig.Lineage.Parent(self.OldTarget).Get()
	if !hasParent {
		return fmt.Errorf("branch %q has no parent", self.Branch)
	}
	args.PrependOpcodes(&ProposalUpdateTarget{
		NewBranch: parent,
		OldBranch: self.OldTarget,
		Proposal:  self.Proposal,
	})
	return nil
}
