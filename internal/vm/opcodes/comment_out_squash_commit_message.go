package opcodes

import (
	"fmt"

	"github.com/git-town/git-town/v16/internal/messages"
	"github.com/git-town/git-town/v16/internal/vm/shared"
)

type CommentOutSquashCommitMessage struct {
	undeclaredOpcodeMethods
}

func (self *CommentOutSquashCommitMessage) Run(args shared.RunArgs) error {
	if err := args.Git.CommentOutSquashCommitMessage(""); err != nil {
		return fmt.Errorf(messages.SquashMessageProblem, err)
	}
	return nil
}
