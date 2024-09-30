package opcodes

import (
	"github.com/git-town/git-town/v16/internal/git"
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/vm/shared"
	. "github.com/git-town/git-town/v16/pkg/prelude"
)

// Commit commits all open changes as a new commit.
// It does not ask the user for a commit message, but chooses one automatically.
type Commit struct {
	Message                 Option[gitdomain.CommitMessage]
	UseDefaultMessage       git.CommitUseDefaultMessage
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *Commit) Run(args shared.RunArgs) error {
	return args.Git.Commit(args.Frontend, self.Message, git.UseDefaultMessageYes, None[gitdomain.Author]())
}
