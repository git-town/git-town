package opcode

import (
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/vm/shared"
)

// IfBranchHasUnmergedChanges allows running different opcodes based on a condition evaluated at runtime.
type IfBranchHasUnmergedChanges struct {
	Branch    domain.LocalBranchName
	WhenTrue  []shared.Opcode // the opcodes to execute if the given branch is empty
	WhenFalse []shared.Opcode // the opcodes to execute if the given branch is not empty
	undeclaredOpcodeMethods
}

func (self *IfBranchHasUnmergedChanges) Run(args shared.RunArgs) error {
	parent := args.Lineage.Parent(self.Branch)
	hasUnmergedChanges, err := args.Runner.Backend.BranchHasUnmergedChanges(self.Branch, parent)
	if err != nil {
		return err
	}
	if hasUnmergedChanges {
		args.PrependOpcodes(self.WhenTrue...)
	} else {
		args.PrependOpcodes(self.WhenFalse...)
	}
	return nil
}
