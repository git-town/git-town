package opcodes

import (
	"github.com/git-town/git-town/v15/internal/git/gitdomain"
	"github.com/git-town/git-town/v15/internal/vm/shared"
	. "github.com/git-town/git-town/v15/pkg/prelude"
)

// CommitOpenChanges commits all open changes as a new commit.
type CommitSquashedChanges struct {
	Message                 Option[gitdomain.CommitMessage]
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *CommitSquashedChanges) Run(args shared.RunArgs) error {
	return args.Git.Commit(args.Frontend, self.Message, "")
}
