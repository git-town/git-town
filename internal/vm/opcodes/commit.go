package opcodes

import (
	"fmt"

	"github.com/git-town/git-town/v23/internal/config/configdomain"
	"github.com/git-town/git-town/v23/internal/git/gitdomain"
	"github.com/git-town/git-town/v23/internal/vm/shared"
	. "github.com/git-town/git-town/v23/pkg/prelude"
)

// Commit commits all open changes as a new commit.
// If no commit message is given, uses FallbackToDefaultCommitMessage to use the default commit message.
//
// If you have a commit message, consider using CommitWithMessage.
type Commit struct {
	AuthorOverride                 Option[gitdomain.Author]
	FallbackToDefaultCommitMessage bool
	Message                        Option[gitdomain.CommitMessage]
}

func (self *Commit) Continue() []shared.Opcode {
	fmt.Println("222222222222222222222222222222222222")

	return []shared.Opcode{CommitIfNeeded{
		AuthorOverride:                 self.AuthorOverride,
		FallbackToDefaultCommitMessage: self.FallbackToDefaultCommitMessage,
		Message:                        self.Message,
	}}
}

func (self *Commit) Run(args shared.RunArgs) error {
	fmt.Println("111111111111111111111111111111111111")
	return args.Git.Commit(args.Frontend, configdomain.UseMessageWithFallbackToDefault(self.Message, self.FallbackToDefaultCommitMessage), self.AuthorOverride, configdomain.CommitHookEnabled)
}
