package opcodes

import (
	"errors"
	"fmt"

	"github.com/git-town/git-town/v14/src/cli/dialog"
	"github.com/git-town/git-town/v14/src/git/gitdomain"
	. "github.com/git-town/git-town/v14/src/gohacks/prelude"
	"github.com/git-town/git-town/v14/src/messages"
	"github.com/git-town/git-town/v14/src/vm/shared"
)

// SquashMerge squash merges the branch with the given name into the current branch.
type SquashMerge struct {
	Branch        gitdomain.LocalBranchName
	CommitMessage Option[gitdomain.CommitMessage]
	Parent        gitdomain.LocalBranchName
	undeclaredOpcodeMethods
}

func (self *SquashMerge) CreateAbortProgram() []shared.Opcode {
	return []shared.Opcode{
		&DiscardOpenChanges{},
	}
}

func (self *SquashMerge) CreateAutomaticUndoError() error {
	return errors.New(messages.ShipAbortedMergeError)
}

func (self *SquashMerge) Run(args shared.RunArgs) error {
	err := args.Git.SquashMerge(args.Frontend, self.Branch)
	if err != nil {
		return err
	}
	branchAuthors, err := args.Git.BranchAuthors(args.Backend, self.Branch, self.Parent)
	if err != nil {
		return err
	}
	author, aborted, err := dialog.SelectSquashCommitAuthor(self.Branch, branchAuthors, args.DialogTestInputs.Value.Next())
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
	if repoAuthor == author {
		author = ""
	}
	err = args.Git.Commit(args.Frontend, self.CommitMessage, author)
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

func (self *SquashMerge) ShouldAutomaticallyUndoOnError() bool {
	return true
}
