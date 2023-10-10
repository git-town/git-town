package step

import (
	"fmt"

	"github.com/git-town/git-town/v9/src/dialog"
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/messages"
)

// SquashMerge squash merges the branch with the given name into the current branch.
type SquashMerge struct {
	Branch        domain.LocalBranchName
	CommitMessage string
	Parent        domain.LocalBranchName
	Empty
}

func (step *SquashMerge) CreateAbortSteps() []Step {
	return []Step{
		&DiscardOpenChanges{},
	}
}

func (step *SquashMerge) CreateAutomaticAbortError() error {
	return fmt.Errorf(messages.ShipAbortedMergeError)
}

func (step *SquashMerge) Run(args RunArgs) error {
	err := args.Runner.Frontend.SquashMerge(step.Branch)
	if err != nil {
		return err
	}
	branchAuthors, err := args.Runner.Backend.BranchAuthors(step.Branch, step.Parent)
	if err != nil {
		return err
	}
	author, err := dialog.SelectSquashCommitAuthor(step.Branch, branchAuthors)
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
	err = args.Runner.Frontend.Commit(step.CommitMessage, author)
	if err != nil {
		return err
	}
	squashedCommitSHA, err := args.Runner.Backend.SHAForBranch(step.Parent.BranchName())
	if err != nil {
		return err
	}
	args.RegisterUndoablePerennialCommit(squashedCommitSHA)
	return nil
}

func (step *SquashMerge) ShouldAutomaticallyAbortOnError() bool {
	return true
}
