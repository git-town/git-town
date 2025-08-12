package opcodes

import (
	"errors"
	"fmt"

	"github.com/git-town/git-town/v21/internal/cli/dialog"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/messages"
	"github.com/git-town/git-town/v21/internal/vm/shared"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

// MergeSquashProgram prepends the opcodes to squash merge the branch with the given name into the current branch.
type MergeSquashProgram struct {
	Authors       []gitdomain.Author
	Branch        gitdomain.LocalBranchName
	CommitMessage Option[gitdomain.CommitMessage]
	Parent        gitdomain.LocalBranchName
	undeclaredOpcodeMethods
}

func (self *MergeSquashProgram) Run(args shared.RunArgs) error {
	author, exit, err := dialog.SquashCommitAuthor(self.Branch, self.Authors, args.Inputs)
	if err != nil {
		return fmt.Errorf(messages.SquashCommitAuthorProblem, err)
	}
	if exit {
		return errors.New("aborted by user")
	}
	repoAuthor := args.Config.Value.ValidatedConfigData.Author()
	var authorOpt Option[gitdomain.Author]
	if repoAuthor == author {
		authorOpt = None[gitdomain.Author]()
	} else {
		authorOpt = Some(author)
	}
	program := []shared.Opcode{
		&MergeSquashAutoUndo{
			Branch: self.Branch,
		},
	}
	if !args.Config.Value.NormalConfig.DryRun {
		program = append(program, &CommitMessageCommentOut{})
	}
	program = append(program,
		&CommitAutoUndo{
			AuthorOverride:                 authorOpt,
			FallbackToDefaultCommitMessage: false,
			Message:                        self.CommitMessage,
		},
		&RegisterUndoablePerennialCommit{
			Parent: self.Parent.BranchName(),
		},
	)
	args.PrependOpcodes(program...)
	return nil
}
