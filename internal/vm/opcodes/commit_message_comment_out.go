package opcodes

import (
	"fmt"

	"github.com/git-town/git-town/v24/internal/messages"
	"github.com/git-town/git-town/v24/internal/vm/shared"
	. "github.com/git-town/git-town/v24/pkg/prelude"
)

// CommitMessageCommentOut comments out the currently active commit message.
type CommitMessageCommentOut struct{}

func (self *CommitMessageCommentOut) Run(args shared.RunArgs) error {
	if err := args.Git.CommentOutSquashCommitMessage(args.Backend, None[string]()); err != nil {
		return fmt.Errorf(messages.SquashMessageProblem, err)
	}
	return nil
}
