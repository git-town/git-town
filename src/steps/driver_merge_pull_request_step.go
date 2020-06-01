package steps

import (
	"fmt"

	"github.com/git-town/git-town/src/drivers"
	"github.com/git-town/git-town/src/git"
)

// DriverMergePullRequestStep squash merges the branch with the given name into the current branch
type DriverMergePullRequestStep struct {
	NoOpStep
	BranchName                string
	PullRequestNumber         int64
	CommitMessage             string
	DefaultCommitMessage      string
	enteredEmptyCommitMessage bool
	mergeError                error
	mergeSha                  string
}

// CreateAbortStep returns the abort step for this step.
func (step *DriverMergePullRequestStep) CreateAbortStep() Step {
	if step.enteredEmptyCommitMessage {
		return &DiscardOpenChangesStep{}
	}
	return nil
}

// CreateUndoStep returns the undo step for this step.
func (step *DriverMergePullRequestStep) CreateUndoStep() Step {
	return &RevertCommitStep{Sha: step.mergeSha}
}

// GetAutomaticAbortErrorMessage returns the error message to display when this step
// cause the command to automatically abort.
func (step *DriverMergePullRequestStep) GetAutomaticAbortErrorMessage() string {
	if step.enteredEmptyCommitMessage {
		return "Aborted because commit exited with error"
	}
	return step.mergeError.Error()
}

// Run executes this step.
func (step *DriverMergePullRequestStep) Run(repo *git.ProdRepo) error {
	commitMessage := step.CommitMessage
	if commitMessage == "" {
		// Allow the user to enter the commit message as if shipping without a driver
		// then revert the commit since merging via the driver will perform the actual squash merge
		step.enteredEmptyCommitMessage = true
		err := repo.Logging.SquashMerge(step.BranchName)
		if err != nil {
			return err
		}
		err = repo.Silent.CommentOutSquashCommitMessage(step.DefaultCommitMessage + "\n\n")
		if err != nil {
			return fmt.Errorf("cannot comment out the squash commit message: %w", err)
		}
		err = repo.Logging.StartCommit()
		if err != nil {
			return err
		}
		commitMessage, err = repo.Silent.LastCommitMessage()
		if err != nil {
			return err
		}
		err = repo.Logging.DeleteLastCommit()
		if err != nil {
			return err
		}
		step.enteredEmptyCommitMessage = false
	}
	driver := drivers.GetActiveDriver()
	currentBranch, err := repo.Silent.CurrentBranch()
	if err != nil {
		return err
	}
	step.mergeSha, step.mergeError = driver.MergePullRequest(drivers.MergePullRequestOptions{
		Branch:            step.BranchName,
		PullRequestNumber: step.PullRequestNumber,
		CommitMessage:     commitMessage,
		LogRequests:       true,
		ParentBranch:      currentBranch,
	})
	return step.mergeError
}

// ShouldAutomaticallyAbortOnError returns whether this step should cause the command to
// automatically abort if it errors.
func (step *DriverMergePullRequestStep) ShouldAutomaticallyAbortOnError() bool {
	return true
}
