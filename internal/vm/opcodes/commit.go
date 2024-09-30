package opcodes

import (
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/vm/shared"
	. "github.com/git-town/git-town/v16/pkg/prelude"
)

// Commit commits all open changes as a new commit.
type Commit struct {
	Message                 Option[gitdomain.CommitMessage]
	UseDefaultCommitMessage gitdomain.UseDefaultCommitMessage
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *Commit) Run(args shared.RunArgs) error {
	return args.Git.Commit(args.Frontend, self.Message, self.UseDefaultCommitMessage, None[gitdomain.Author]())
}
