package opcodes

import (
	"errors"
	"fmt"

	"github.com/git-town/git-town/v12/src/cli/dialog"
	"github.com/git-town/git-town/v12/src/git/gitdomain"
	"github.com/git-town/git-town/v12/src/messages"
	"github.com/git-town/git-town/v12/src/vm/shared"
)

// SquashMerge squash merges the branch with the given name into the current branch.
type SquashMerge struct {
	Branch        gitdomain.LocalBranchName
	CommitMessage string
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
	err := args.Runner.Frontend.SquashMerge(self.Branch)
	if err != nil {
		return err
	}
	branchAuthors, err := args.Runner.Backend.BranchAuthors(self.Branch, self.Parent)
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
	repoAuthor, err := args.Runner.Backend.Author()
	if err != nil {
		return err
	}
	if !args.Runner.DryRun {
		if err = args.Runner.Backend.CommentOutSquashCommitMessage(""); err != nil {
			return fmt.Errorf(messages.SquashMessageProblem, err)
		}
	}
	if repoAuthor == author {
		author = ""
	}
	err = args.Runner.Frontend.Commit(self.CommitMessage, author)
	if err != nil {
		return err
	}
	squashedCommitSHA, err := args.Runner.Backend.SHAForBranch(self.Parent.BranchName())
	if err != nil {
		return err
	}
	args.RegisterUndoablePerennialCommit(squashedCommitSHA)
	return nil
}

func (self *SquashMerge) ShouldAutomaticallyUndoOnError() bool {
	return true
}
