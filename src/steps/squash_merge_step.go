package steps

import (
	"fmt"

	"github.com/git-town/git-town/v9/src/dialog"
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/messages"
)

// SquashMergeStep squash merges the branch with the given name into the current branch.
type SquashMergeStep struct {
	Branch        domain.LocalBranchName
	CommitMessage string
	Parent        domain.LocalBranchName
	EmptyStep
}

func (step *SquashMergeStep) CreateAbortSteps() []Step {
	return []Step{&DiscardOpenChangesStep{}}
}

func (step *SquashMergeStep) CreateAutomaticAbortError() error {
	return fmt.Errorf(messages.ShipAbortedMergeError)
}

func (step *SquashMergeStep) Run(args RunArgs) error {
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
	// TODO: read this SHA from the output of commit above and get rid of the SHAForBranch call below
	squashedCommitSHA, err := args.Runner.Backend.SHAForBranch(step.Parent.BranchName())
	if err != nil {
		return err
	}
	args.RegisterUndoablePerennialCommit(squashedCommitSHA)
	return nil
}

func (step *SquashMergeStep) ShouldAutomaticallyAbortOnError() bool {
	return true
}
