package opcodes

import (
	"fmt"

	"github.com/git-town/git-town/v16/internal/messages"
	"github.com/git-town/git-town/v16/internal/vm/shared"
)

// Commit commits all open changes as a new commit.
// If no commit message is given, uses FallbackToDefaultCommitMessage to use the default commit message.
//
// If you have a commit message, consider using CommitWithMessage.
type CommitMessageCommentOut struct {
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *CommitMessageCommentOut) Run(args shared.RunArgs) error {
	if err := args.Git.CommentOutSquashCommitMessage(""); err != nil {
		return fmt.Errorf(messages.SquashMessageProblem, err)
	}
	return nil
}
