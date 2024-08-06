package opcodes

import (
	"github.com/git-town/git-town/v14/internal/git/gitdomain"
	"github.com/git-town/git-town/v14/internal/vm/shared"
)

// Merge merges the branch with the given name into the current branch.
type Merge struct {
	Branch                  gitdomain.BranchName
	undeclaredOpcodeMethods `exhaustruct:"optional"`
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
	return args.Git.MergeBranchNoEdit(args.Frontend, self.Branch)
}
