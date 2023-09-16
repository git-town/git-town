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
	return args.Runner.Frontend.Commit(step.CommitMessage, author)
}

func (step *SquashMergeStep) ShouldAutomaticallyAbortOnError() bool {
	return true
}
