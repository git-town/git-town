package opcode

import (
	"fmt"

	"github.com/git-town/git-town/v9/src/dialog"
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/messages"
	"github.com/git-town/git-town/v9/src/vm/shared"
)

// SquashMerge squash merges the branch with the given name into the current branch.
type SquashMerge struct {
	Branch        domain.LocalBranchName
	CommitMessage string
	Parent        domain.LocalBranchName
	undeclaredOpcodeMethods
}

func (op *SquashMerge) CreateAbortProgram() []shared.Opcode {
	return []shared.Opcode{
		&DiscardOpenChanges{},
	}
}

func (op *SquashMerge) CreateAutomaticAbortError() error {
	return fmt.Errorf(messages.ShipAbortedMergeError)
}

func (op *SquashMerge) Run(args shared.RunArgs) error {
	err := args.Runner.Frontend.SquashMerge(op.Branch)
	if err != nil {
		return err
	}
	branchAuthors, err := args.Runner.Backend.BranchAuthors(op.Branch, op.Parent)
	if err != nil {
		return err
	}
	author, err := dialog.SelectSquashCommitAuthor(op.Branch, branchAuthors)
	if err != nil {
		return fmt.Errorf(messages.SquashCommitAuthorProblem, err)
	}
	repoAuthor, err := args.Runner.Backend.Author()
	if err != nil {
		return fmt.Errorf(messages.GitUserProblem, err)
	}
	if err = args.Runner.Backend.CommentOutSquashCommitMessage(""); err != nil {
		return fmt.Errorf(messages.SquashMessageProblem, err)
	}
	if repoAuthor == author {
		author = ""
	}
	err = args.Runner.Frontend.Commit(op.CommitMessage, author)
	if err != nil {
		return err
	}
	squashedCommitSHA, err := args.Runner.Backend.SHAForBranch(op.Parent.BranchName())
	if err != nil {
		return err
	}
	args.RegisterUndoablePerennialCommit(squashedCommitSHA)
	return nil
}

func (op *SquashMerge) ShouldAutomaticallyAbortOnError() bool {
	return true
}
