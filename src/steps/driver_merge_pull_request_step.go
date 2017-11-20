package steps

import (
	"github.com/Originate/exit"
	"github.com/Originate/git-town/src/drivers"
	"github.com/Originate/git-town/src/git"
	"github.com/Originate/git-town/src/script"
)

// DriverMergePullRequestStep squash merges the branch with the given name into the current branch
type DriverMergePullRequestStep struct {
	NoOpStep
	BranchName                string
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

// CreateUndoStepAfterRun returns the undo step for this step after it is run.
func (step *DriverMergePullRequestStep) CreateUndoStepAfterRun() Step {
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
func (step *DriverMergePullRequestStep) Run() error {
	commitMessage := step.CommitMessage
	if commitMessage == "" {
		// Allow the user to enter the commit message as if shipping without a driver
		// then revert the commit since merging via the driver will perform the actual squash merge
		step.enteredEmptyCommitMessage = true
		script.SquashMerge(step.BranchName)
		git.CommentOutSquashCommitMessage(step.DefaultCommitMessage + "\n\n")
		err := script.RunCommand("git", "commit")
		if err != nil {
			return err
		}
		commitMessage = git.GetLastCommitMessage()
		err = script.RunCommand("git", "reset", "--hard", "HEAD~1")
		exit.IfWrap(err, "Error resetting the main branch")
		step.enteredEmptyCommitMessage = false
	}
	driver := drivers.GetActiveDriver()
	step.mergeSha, step.mergeError = driver.MergePullRequest(drivers.MergePullRequestOptions{
		Branch:        step.BranchName,
		CommitMessage: commitMessage,
		LogRequests:   true,
		ParentBranch:  git.GetCurrentBranchName(),
	})
	return step.mergeError
}

// ShouldAutomaticallyAbortOnError returns whether this step should cause the command to
// automatically abort if it errors.
func (step *DriverMergePullRequestStep) ShouldAutomaticallyAbortOnError() bool {
	return true
}
