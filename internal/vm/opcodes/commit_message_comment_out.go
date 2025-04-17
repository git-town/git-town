package opcodes

import (
	"fmt"

	"github.com/git-town/git-town/v19/internal/messages"
	"github.com/git-town/git-town/v19/internal/vm/shared"
	. "github.com/git-town/git-town/v19/pkg/prelude"
)

// comments out the currently active commit message
type CommitMessageCommentOut struct {
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *CommitMessageCommentOut) Run(args shared.RunArgs) error {
	if err := args.Git.CommentOutSquashCommitMessage(None[string]()); err != nil {
		return fmt.Errorf(messages.SquashMessageProblem, err)
	}
	return nil
}
