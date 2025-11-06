package opcodes

import (
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/vm/shared"
)

// CommitRemove removes the commit with the given SHA from the given branch.
type CommitRemove struct {
	SHA gitdomain.SHA
}

func (self *CommitRemove) Abort() []shared.Opcode {
	return []shared.Opcode{
		&RebaseAbort{},
	}
}

func (self *CommitRemove) Continue() []shared.Opcode {
	return []shared.Opcode{
		&RebaseContinueIfNeeded{},
	}
}

func (self *CommitRemove) Run(args shared.RunArgs) error {
	return args.Git.RemoveCommit(args.Frontend, self.SHA)
}
