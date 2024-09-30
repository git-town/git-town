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

// SquashMerge squash merges the branch with the given name into the current branch.
type SquashMerge struct {
	Branch        gitdomain.LocalBranchName
	CommitMessage Option[gitdomain.CommitMessage]
	Parent        gitdomain.LocalBranchName
	undeclaredOpcodeMethods
}

func (self *SquashMerge) AbortProgram() []shared.Opcode {
	return []shared.Opcode{
		&DiscardOpenChanges{},
	}
}

func (self *SquashMerge) AutomaticUndoError() error {
	return errors.New(messages.ShipAbortedMergeError)
}

func (self *SquashMerge) Run(args shared.RunArgs) error {
	// TODO: extract into separate opcodes for Git resilience
	// Possible create a SquashMergeProgram function that returns these opcodes
	err := args.Git.SquashMerge(args.Frontend, self.Branch)
	if err != nil {
		return err
	}
	branchAuthors, err := args.Git.BranchAuthors(args.Backend, self.Branch, self.Parent)
	if err != nil {
		return err
	}
	author, aborted, err := dialog.SelectSquashCommitAuthor(self.Branch, branchAuthors, args.DialogTestInputs.Next())
	if err != nil {
		return fmt.Errorf(messages.SquashCommitAuthorProblem, err)
	}
	if aborted {
		return errors.New("aborted by user")
	}
	repoAuthor := args.Config.Author()
	if !args.Config.DryRun {
		if err = args.Git.CommentOutSquashCommitMessage(""); err != nil {
			return fmt.Errorf(messages.SquashMessageProblem, err)
		}
	}
	var authorOpt Option[gitdomain.Author]
	if repoAuthor == author {
		authorOpt = None[gitdomain.Author]()
	} else {
		authorOpt = Some(author)
	}
	err = args.Git.Commit(args.Frontend, self.CommitMessage, gitdomain.FallbackToDefaultCommitMessageNo, authorOpt)
	if err != nil {
		return err
	}
	squashedCommitSHA, err := args.Git.SHAForBranch(args.Backend, self.Parent.BranchName())
	if err != nil {
		return err
	}
	args.RegisterUndoablePerennialCommit(squashedCommitSHA)
	return nil
}

func (self *SquashMerge) ShouldUndoOnError() bool {
	return true
}
