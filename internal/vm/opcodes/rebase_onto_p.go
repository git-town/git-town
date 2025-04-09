package opcodes

import (
	"github.com/git-town/git-town/v18/internal/git/gitdomain"
	"github.com/git-town/git-town/v18/internal/vm/shared"
	. "github.com/git-town/git-town/v18/pkg/prelude"
)

// rebases the current branch against the target branch while executing "git town swap", while moving the target branch onto the Onto branch.
type RebaseOntoP struct {
	BranchToRebaseOnto      gitdomain.Location
	CommitsToRemove         gitdomain.Location
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *RebaseOntoP) AbortProgram() []shared.Opcode {
	return []shared.Opcode{
		&RebaseAbort{},
	}
}

func (self *RebaseOntoP) ContinueProgram() []shared.Opcode {
	return []shared.Opcode{
		&RebaseContinueIfNeeded{},
	}
}

func (self *RebaseOntoP) Run(args shared.RunArgs) error {
	return args.Git.RebaseOntoP(args.Frontend, self.BranchToRebaseOnto, self.CommitsToRemove, None[gitdomain.LocalBranchName]())
}
