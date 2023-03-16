package steps

import (
	"fmt"

	"github.com/git-town/git-town/v7/src/dialog"
	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/hosting"
)

// SquashMergeStep squash merges the branch with the given name into the current branch.
type SquashMergeStep struct {
	EmptyStep
	Branch        string
	CommitMessage string
	Parent        string
}

func (step *SquashMergeStep) CreateAbortStep() Step {
	return &DiscardOpenChangesStep{}
}

func (step *SquashMergeStep) CreateUndoStep(repo *git.PublicRepo) (Step, error) {
	currentSHA, err := repo.CurrentSha()
	if err != nil {
		return nil, err
	}
	return &RevertCommitStep{Sha: currentSHA}, nil
}

func (step *SquashMergeStep) CreateAutomaticAbortError() error {
	return fmt.Errorf("aborted because commit exited with error")
}

func (step *SquashMergeStep) Run(repo *git.PublicRepo, connector hosting.Connector) error {
	err := repo.SquashMerge(step.Branch)
	if err != nil {
		return err
	}
	branchAuthors, err := repo.BranchAuthors(step.Branch, step.Parent)
	if err != nil {
		return err
	}
	author, err := dialog.SelectSquashCommitAuthor(step.Branch, branchAuthors)
	if err != nil {
		return fmt.Errorf("error getting squash commit author: %w", err)
	}
	repoAuthor, err := repo.Author()
	if err != nil {
		return fmt.Errorf("cannot determine repo author: %w", err)
	}
	if err = repo.CommentOutSquashCommitMessage(""); err != nil {
		return fmt.Errorf("cannot comment out the squash commit message: %w", err)
	}
	if repoAuthor == author {
		author = ""
	}
	return repo.Commit(step.CommitMessage, author)
}

func (step *SquashMergeStep) ShouldAutomaticallyAbortOnError() bool {
	return true
}
