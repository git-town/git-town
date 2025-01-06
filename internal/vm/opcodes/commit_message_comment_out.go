package opcodes

import (
	"fmt"

	"github.com/git-town/git-town/v17/internal/messages"
	"github.com/git-town/git-town/v17/internal/vm/shared"
)

// comments out the currently active commit message
type CommitMessageCommentOut struct {
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *CommitMessageCommentOut) Run(args shared.RunArgs) error {
	if err := args.Git.CommentOutSquashCommitMessage(""); err != nil {
		return fmt.Errorf(messages.SquashMessageProblem, err)
	}
	return nil
}
