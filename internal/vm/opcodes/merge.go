package opcodes

import (
	"github.com/git-town/git-town/v17/internal/git/gitdomain"
	"github.com/git-town/git-town/v17/internal/vm/shared"
)

// Merge merges the branch with the given name into the current branch.
type Merge struct {
	Branch                  gitdomain.BranchName
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *Merge) AbortProgram() []shared.Opcode {
	return []shared.Opcode{
		&MergeAbort{},
	}
}

func (self *Merge) ContinueProgram() []shared.Opcode {
	return []shared.Opcode{
		&MergeContinue{},
	}
}

func (self *Merge) Run(args shared.RunArgs) error {
	return args.Git.MergeBranchNoEdit(args.Frontend, self.Branch)
}
