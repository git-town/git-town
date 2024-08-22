package opcodes

import (
	"github.com/git-town/git-town/v15/internal/git/gitdomain"
	"github.com/git-town/git-town/v15/internal/vm/shared"
)

// CommitOpenChanges commits all open changes as a new commit.
// It does not ask the user for a commit message, but chooses one automatically.
type CommitOpenChanges struct {
	Message                 gitdomain.CommitMessage
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *CommitOpenChanges) Run(args shared.RunArgs) error {
	return args.Git.CommitStagedChanges(args.Frontend, self.Message)
}
