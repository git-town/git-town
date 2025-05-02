package opcodes

import (
	"github.com/git-town/git-town/v20/internal/config/configdomain"
	"github.com/git-town/git-town/v20/internal/git/gitdomain"
	"github.com/git-town/git-town/v20/internal/vm/shared"
	. "github.com/git-town/git-town/v20/pkg/prelude"
)

// Commit commits all open changes as a new commit.
// If no commit message is given, uses FallbackToDefaultCommitMessage to use the default commit message.
//
// If you have a commit message, consider using CommitWithMessage.
type Commit struct {
	AuthorOverride                 Option[gitdomain.Author]
	FallbackToDefaultCommitMessage bool
	Message                        Option[gitdomain.CommitMessage]
	undeclaredOpcodeMethods        `exhaustruct:"optional"`
}

func (self *Commit) Run(args shared.RunArgs) error {
	return args.Git.Commit(args.Frontend, configdomain.UseMessageWithFallbackToDefault(self.Message, self.FallbackToDefaultCommitMessage), self.AuthorOverride, configdomain.CommitHookEnabled)
}
