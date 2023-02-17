package steps

import (
	"fmt"

	"github.com/git-town/git-town/v7/src/dialog"
	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/hosting"
)

// SquashMergeBranchStep squash merges the branch with the given name into the current branch.
type SquashMergeBranchStep struct {
	NoOpStep
	Branch        string
	CommitMessage string
}

func (step *SquashMergeBranchStep) CreateAbortStep() Step { //nolint:ireturn
	return &DiscardOpenChangesStep{}
}

func (step *SquashMergeBranchStep) CreateUndoStep(repo *git.ProdRepo) (Step, error) { //nolint:ireturn
	currentSHA, err := repo.Silent.CurrentSha()
	if err != nil {
		return nil, err
	}
	return &RevertCommitStep{Sha: currentSHA}, nil
}

func (step *SquashMergeBranchStep) CreateAutomaticAbortError() error {
	return fmt.Errorf("aborted because commit exited with error")
}

func (step *SquashMergeBranchStep) Run(repo *git.ProdRepo, driver hosting.Driver) error {
	err := repo.Logging.SquashMerge(step.Branch)
	if err != nil {
		return err
	}
	author, err := dialog.DetermineSquashCommitAuthor(step.Branch, repo)
	if err != nil {
		return fmt.Errorf("error getting squash commit author: %w", err)
	}
	repoAuthor, err := repo.Silent.Author()
	if err != nil {
		return fmt.Errorf("cannot determine repo author: %w", err)
	}
	if err = repo.Silent.CommentOutSquashCommitMessage(""); err != nil {
		return fmt.Errorf("cannot comment out the squash commit message: %w", err)
	}
	if repoAuthor == author {
		author = ""
	}
	return repo.Logging.Commit(step.CommitMessage, author)
}

func (step *SquashMergeBranchStep) ShouldAutomaticallyAbortOnError() bool {
	return true
}
