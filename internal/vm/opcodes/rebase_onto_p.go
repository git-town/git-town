package opcodes

import (
	"github.com/git-town/git-town/v18/internal/git/gitdomain"
	"github.com/git-town/git-town/v18/internal/vm/shared"
)

// rebases the current branch against the target branch while executing "git town swap", while moving the target branch onto the Onto branch.
type RemoveCommit struct {
	Commit                  gitdomain.SHA
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *RemoveCommit) AbortProgram() []shared.Opcode {
	return []shared.Opcode{
		&RebaseAbort{},
	}
}

func (self *RemoveCommit) ContinueProgram() []shared.Opcode {
	return []shared.Opcode{
		&RebaseContinueIfNeeded{},
	}
}

func (self *RemoveCommit) Run(args shared.RunArgs) error {
	return args.Git.RemoveCommit(args.Frontend, self.Commit)
}
