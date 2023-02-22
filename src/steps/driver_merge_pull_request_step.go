package steps

import (
	"fmt"

	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/hosting"
)

// DriverSquashMergeProposalStep squash merges the branch with the given name into the current branch.
type DriverSquashMergeProposalStep struct {
	NoOpStep
	Branch                    string
	CommitMessage             string
	DefaultProposalMessage    string
	enteredEmptyCommitMessage bool
	mergeError                error
	mergeSha                  string
	PullRequestNumber         int
}

func (step *DriverSquashMergeProposalStep) CreateAbortStep() Step {
	if step.enteredEmptyCommitMessage {
		return &DiscardOpenChangesStep{}
	}
	return nil
}

func (step *DriverSquashMergeProposalStep) CreateUndoStep(repo *git.ProdRepo) (Step, error) {
	return &RevertCommitStep{Sha: step.mergeSha}, nil
}

func (step *DriverSquashMergeProposalStep) CreateAutomaticAbortError() error {
	if step.enteredEmptyCommitMessage {
		return fmt.Errorf("aborted because commit exited with error")
	}
	return step.mergeError
}

func (step *DriverSquashMergeProposalStep) Run(repo *git.ProdRepo, driver hosting.Driver) error {
	commitMessage := step.CommitMessage
	//nolint:nestif
	if commitMessage == "" {
		// Allow the user to enter the commit message as if shipping without a connector
		// then revert the commit since merging via the connector will perform the actual squash merge.
		step.enteredEmptyCommitMessage = true
		err := repo.Logging.SquashMerge(step.Branch)
		if err != nil {
			return err
		}
		err = repo.Silent.CommentOutSquashCommitMessage(step.DefaultProposalMessage + "\n\n")
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
	currentBranch, err := repo.Silent.CurrentBranch()
	if err != nil {
		return err
	}
	step.mergeSha, step.mergeError = driver.SquashMergeProposal(hosting.SquashMergeProposalOptions{
		Branch:            step.Branch,
		PullRequestNumber: step.PullRequestNumber,
		CommitMessage:     commitMessage,
		LogRequests:       true,
		ParentBranch:      currentBranch,
	})
	return step.mergeError
}

// ShouldAutomaticallyAbortOnError returns whether this step should cause the command to
// automatically abort if it errors.
func (step *DriverSquashMergeProposalStep) ShouldAutomaticallyAbortOnError() bool {
	return true
}
