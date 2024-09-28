package opcodes

import (
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/vm/shared"
	. "github.com/git-town/git-town/v16/pkg/prelude"
)

// CommitOpenChanges commits all open changes as a new commit.
// It does not ask the user for a commit message, but chooses one automatically.
type CommitOpenChanges struct {
	Message                 Option[gitdomain.CommitMessage]
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *CommitOpenChanges) Run(args shared.RunArgs) error {
	return args.Git.CommitNeverAskUser(args.Frontend, self.Message)
}
