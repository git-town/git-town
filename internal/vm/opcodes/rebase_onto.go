package opcodes

import (
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/vm/shared"
)

// rebases the current branch against the target branch, while moving the target branch onto the Onto branch.
type RebaseOnto struct {
	BranchToRebaseAgainst   gitdomain.BranchName
	BranchToRebaseOnto      gitdomain.LocalBranchName
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *RebaseOnto) AbortProgram() []shared.Opcode {
	return []shared.Opcode{
		&RebaseAbort{},
	}
}

func (self *RebaseOnto) ContinueProgram() []shared.Opcode {
	return []shared.Opcode{
		&RebaseContinueIfNeeded{},
	}
}

func (self *RebaseOnto) Run(args shared.RunArgs) error {
	return args.Git.RebaseOnto(args.Frontend, self.BranchToRebaseAgainst, self.BranchToRebaseOnto)
}
