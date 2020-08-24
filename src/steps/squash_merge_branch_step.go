package steps

import (
	"fmt"

	"github.com/git-town/git-town/src/drivers"
	"github.com/git-town/git-town/src/git"
	"github.com/git-town/git-town/src/prompt"
)

// SquashMergeBranchStep squash merges the branch with the given name into the current branch.
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
func (step *SquashMergeBranchStep) CreateUndoStep(repo *git.ProdRepo) (Step, error) {
	currentSHA, err := repo.Silent.CurrentSha()
	if err != nil {
		return nil, err
	}
	return &RevertCommitStep{Sha: currentSHA}, nil
}

// GetAutomaticAbortError returns the error message to display when this step
// cause the command to automatically abort.
func (step *SquashMergeBranchStep) GetAutomaticAbortError() error {
	return fmt.Errorf("aborted because commit exited with error")
}

// Run executes this step.
func (step *SquashMergeBranchStep) Run(repo *git.ProdRepo, driver drivers.CodeHostingDriver) error {
	err := repo.Logging.SquashMerge(step.BranchName)
	if err != nil {
		return err
	}
	author, err := prompt.GetSquashCommitAuthor(step.BranchName, repo)
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
	switch {
	case author != repoAuthor && step.CommitMessage != "":
		return repo.Logging.CommitWithMessageAndAuthor(step.CommitMessage, author)
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
