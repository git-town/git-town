package steps

import (
	"log"
	"strings"

	"github.com/Originate/git-town/src/drivers"
	"github.com/Originate/git-town/src/git"
	"github.com/Originate/git-town/src/script"
	"github.com/Originate/git-town/src/util"
)

// DriverMergePullRequestStep squash merges the branch with the given name into the current branch
type DriverMergePullRequestStep struct {
	NoOpStep
	BranchName                string
	CommitMessage             string
	Driver                    drivers.CodeHostingDriver
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
		step.enteredEmptyCommitMessage = true
		err := script.RunCommand("git", "merge", "--squash", step.BranchName)
		if err != nil {
			log.Fatal("Error squash merging:", err)
		}
		util.GetCommandOutput("sed", "-i", "-e", "s/^/# /g", ".git/SQUASH_MSG")
		err = script.RunCommand("git", "commit")
		if err != nil {
			return err
		}
		commitMessage = util.GetCommandOutput("git", "log", "-1", "--format=%B")
		err = script.RunCommand("git", "reset", "--hard", "HEAD~1")
		if err != nil {
			log.Fatal("Error resetting the main branch", err)
		}
		step.enteredEmptyCommitMessage = false
	}
	repository := git.GetURLRepositoryName(git.GetRemoteOriginURL())
	repositoryParts := strings.SplitN(repository, "/", 2)
	step.mergeSha, step.mergeError = step.Driver.MergePullRequest(drivers.MergePullRequestOptions{
		Branch:        step.BranchName,
		CommitMessage: commitMessage,
		LogRequests:   true,
		Owner:         repositoryParts[0],
		ParentBranch:  git.GetCurrentBranchName(),
		Repository:    repositoryParts[1],
	})
	return step.mergeError
}

// ShouldAutomaticallyAbortOnError returns whether this step should cause the command to
// automatically abort if it errors.
func (step *DriverMergePullRequestStep) ShouldAutomaticallyAbortOnError() bool {
	return true
}
