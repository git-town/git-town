package opcodes

import (
	"github.com/git-town/git-town/v13/src/git/gitdomain"
	"github.com/git-town/git-town/v13/src/vm/shared"
)

// Merge merges the branch with the given name into the current branch.
type Merge struct {
	Branch gitdomain.BranchName
	undeclaredOpcodeMethods
}

func (self *Merge) CreateAbortProgram() []shared.Opcode {
	return []shared.Opcode{
		&AbortMerge{},
	}
}

func (self *Merge) CreateContinueProgram() []shared.Opcode {
	return []shared.Opcode{
		&ContinueMerge{},
	}
}

func (self *Merge) Run(args shared.RunArgs) error {
	return args.Runner.Frontend.MergeBranchNoEdit(self.Branch)
}
