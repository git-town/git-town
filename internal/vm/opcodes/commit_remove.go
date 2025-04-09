package opcodes

import (
	"github.com/git-town/git-town/v18/internal/git/gitdomain"
	"github.com/git-town/git-town/v18/internal/vm/shared"
)

// removes the given commit from the given branch
type CommitRemove struct {
	SHA                     gitdomain.SHA
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *CommitRemove) AbortProgram() []shared.Opcode {
	return []shared.Opcode{
		&RebaseAbort{},
	}
}

func (self *CommitRemove) ContinueProgram() []shared.Opcode {
	return []shared.Opcode{
		&RebaseContinueIfNeeded{},
	}
}

func (self *CommitRemove) Run(args shared.RunArgs) error {
	return args.Git.RemoveCommit(args.Frontend, self.SHA)
}
