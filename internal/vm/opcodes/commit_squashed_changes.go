package opcodes

import (
	"github.com/git-town/git-town/v14/internal/git/gitdomain"
	"github.com/git-town/git-town/v14/internal/vm/shared"
	. "github.com/git-town/git-town/v14/pkg/prelude"
)

// CommitOpenChanges commits all open changes as a new commit.
type CommitSquashedChanges struct {
	Message                 Option[gitdomain.CommitMessage]
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *CommitSquashedChanges) Run(args shared.RunArgs) error {
	return args.Git.Commit(args.Frontend, self.Message, "")
}
