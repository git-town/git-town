package opcodes

import (
	"fmt"

	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/vm/shared"
)

// ProposalUpdateBaseToParent updates the target of the proposal with the given number to the parent of the given branch.
type ProposalUpdateBaseToParent struct {
	Branch                  gitdomain.LocalBranchName
	OldTarget               gitdomain.LocalBranchName
	ProposalNumber          int
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *ProposalUpdateBaseToParent) Run(args shared.RunArgs) error {
	parent, hasParent := args.Config.Value.NormalConfig.Lineage.Parent(self.Branch).Get()
	if !hasParent {
		return fmt.Errorf("branch %q has no parent", self.Branch)
	}
	args.PrependOpcodes(&ProposalUpdateBase{
		NewTarget:      parent,
		OldTarget:      self.OldTarget,
		ProposalNumber: self.ProposalNumber,
	})
	return nil
}
