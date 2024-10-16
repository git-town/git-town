package opcodes

import (
	"errors"
	"fmt"

	"github.com/git-town/git-town/v16/internal/cli/dialog"
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/messages"
	"github.com/git-town/git-town/v16/internal/vm/shared"
	. "github.com/git-town/git-town/v16/pkg/prelude"
)

// MergeSquash squash merges the branch with the given name into the current branch.
type MergeSquash struct {
	Authors       []gitdomain.Author
	Branch        gitdomain.LocalBranchName
	CommitMessage Option[gitdomain.CommitMessage]
	Parent        gitdomain.LocalBranchName
	undeclaredOpcodeMethods
}

func (self *MergeSquash) AbortProgram() []shared.Opcode {
	return []shared.Opcode{
		&ChangesDiscard{},
	}
}

func (self *MergeSquash) AutomaticUndoError() error {
	return errors.New(messages.ShipAbortedMergeError)
}

func (self *MergeSquash) Run(args shared.RunArgs) error {
	// TODO: extract into separate opcodes for Git resilience
	// Possible create a SquashMergeProgram function that returns these opcodes
	author, aborted, err := dialog.SelectSquashCommitAuthor(self.Branch, self.Authors, args.DialogTestInputs.Next())
	if err != nil {
		return fmt.Errorf(messages.SquashCommitAuthorProblem, err)
	}
	if aborted {
		return errors.New("aborted by user")
	}
	repoAuthor := args.Config.Author()
	var authorOpt Option[gitdomain.Author]
	if repoAuthor == author {
		authorOpt = None[gitdomain.Author]()
	} else {
		authorOpt = Some(author)
	}
	err = args.Git.SquashMerge(args.Frontend, self.Branch)
	if err != nil {
		return err
	}
	if !args.Config.DryRun {
		if err = args.Git.CommentOutSquashCommitMessage(""); err != nil {
			return fmt.Errorf(messages.SquashMessageProblem, err)
		}
	}
	args.PrependOpcodes(
		&CommitAutoUndo{
			AuthorOverride:                 authorOpt,
			FallbackToDefaultCommitMessage: false,
			Message:                        self.CommitMessage,
		},
		&RegisterUndoablePerennialCommit{
			Parent: self.Parent.BranchName(),
		})
	return nil
}

func (self *MergeSquash) ShouldUndoOnError() bool {
	return true
}
