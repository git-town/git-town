package steps

import (
	"fmt"

	"github.com/git-town/git-town/src/git"
	"github.com/git-town/git-town/src/prompt"
)

// SquashMergeBranchStep squash merges the branch with the given name into the current branch
type SquashMergeBranchStep struct {
	NoOpStep
	BranchName    string
	CommitMessage string
}

// CreateAbortStep returns the abort step for this step.
func (step *SquashMergeBranchStep) CreateAbortStep() Step {
	return &DiscardOpenChangesStep{}
}

// CreateUndoStep returns the undo step for this step.
func (step *SquashMergeBranchStep) CreateUndoStep() Step {
	return &RevertCommitStep{Sha: git.GetCurrentSha()}
}

// GetAutomaticAbortErrorMessage returns the error message to display when this step
// cause the command to automatically abort.
func (step *SquashMergeBranchStep) GetAutomaticAbortErrorMessage() string {
	return "Aborted because commit exited with error"
}

// Run executes this step.
func (step *SquashMergeBranchStep) Run(repo *git.ProdRepo) error {
	err := repo.Logging.SquashMerge(step.BranchName)
	if err != nil {
		return err
	}
	author := prompt.GetSquashCommitAuthor(step.BranchName)
	repoAuthor, err := repo.Silent.Author()
	if err != nil {
		return err
	}
	if err = repo.Silent.CommentOutSquashCommitMessage(""); err != nil {
		return fmt.Errorf("cannot comment out the squash commit message: %w", err)
	}
	switch {
	case author != repoAuthor && step.CommitMessage != "":
		return repo.Logging.CommitWithMessageAndAuthor(step.CommitMessage, author)
	case author != repoAuthor:
		return repo.Logging.CommitWithAuthor(author)
	case step.CommitMessage != "":
		return repo.Logging.CommitWithMessage(step.CommitMessage)
	default:
		return repo.Logging.Commit()
	}
}

// ShouldAutomaticallyAbortOnError returns whether this step should cause the command to
// automatically abort if it errors.
func (step *SquashMergeBranchStep) ShouldAutomaticallyAbortOnError() bool {
	return true
}
