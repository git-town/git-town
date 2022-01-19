//nolint:ireturn
package steps

import (
	"fmt"

	"github.com/git-town/git-town/v7/src/drivers"
	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/prompt"
)

// SquashMergeBranchStep squash merges the branch with the given name into the current branch.
type SquashMergeBranchStep struct {
	NoOpStep
	BranchName    string
	CommitMessage string
}

func (step *SquashMergeBranchStep) CreateAbortStep() Step {
	return &DiscardOpenChangesStep{}
}

func (step *SquashMergeBranchStep) CreateUndoStep(repo *git.ProdRepo) (Step, error) {
	currentSHA, err := repo.Silent.CurrentSha()
	if err != nil {
		return nil, err
	}
	return &RevertCommitStep{Sha: currentSHA}, nil
}

func (step *SquashMergeBranchStep) CreateAutomaticAbortError() error {
	return fmt.Errorf("aborted because commit exited with error")
}

func (step *SquashMergeBranchStep) Run(repo *git.ProdRepo, driver drivers.CodeHostingDriver) error {
	err := repo.Logging.SquashMerge(step.BranchName)
	if err != nil {
		return err
	}
	author, err := prompt.DetermineSquashCommitAuthor(step.BranchName, repo)
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

func (step *SquashMergeBranchStep) ShouldAutomaticallyAbortOnError() bool {
	return true
}
