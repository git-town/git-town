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

// MergeSquashProgram squash merges the branch with the given name into the current branch.
type MergeSquashProgram struct {
	Authors       []gitdomain.Author
	Branch        gitdomain.LocalBranchName
	CommitMessage Option[gitdomain.CommitMessage]
	Parent        gitdomain.LocalBranchName
	undeclaredOpcodeMethods
}

func (self *MergeSquashProgram) AbortProgram() []shared.Opcode {
	return []shared.Opcode{
		&ChangesDiscard{},
	}
}

func (self *MergeSquashProgram) AutomaticUndoError() error {
	return errors.New(messages.ShipAbortedMergeError)
}

func (self *MergeSquashProgram) Run(args shared.RunArgs) error {
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
	program := []shared.Opcode{
		&MergeSquashAutoUndo{
			Branch: self.Branch,
		},
	}
	if !args.Config.DryRun {
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

func (self *MergeSquashProgram) ShouldUndoOnError() bool {
	return true
}
