package steps

import (
	"github.com/Originate/git-town/src/git"
	"github.com/Originate/git-town/src/prompt"
	"github.com/Originate/git-town/src/script"
	"github.com/pkg/errors"
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

// CreateUndoStepAfterRun returns the undo step for this step after it is run.
func (step *SquashMergeBranchStep) CreateUndoStepAfterRun() Step {
	return &RevertCommitStep{Sha: git.GetCurrentSha()}
}

// GetAutomaticAbortErrorMessage returns the error message to display when this step
// cause the command to automatically abort.
func (step *SquashMergeBranchStep) GetAutomaticAbortErrorMessage() string {
	return "Aborted because commit exited with error"
}

// Run executes this step.
func (step *SquashMergeBranchStep) Run() error {
	err := script.SquashMerge(step.BranchName)
	if err != nil {
		return errors.Wrapf(err, "cannot squash-merge branch %q", step.BranchName)
	}
	args := []string{"commit"}
	if step.CommitMessage != "" {
		args = append(args, "-m", step.CommitMessage)
	}
	author := prompt.GetSquashCommitAuthor(step.BranchName)
	if author != git.GetLocalAuthor() {
		args = append(args, "--author", author)
	}
	err = git.CommentOutSquashCommitMessage("")
	if err != nil {
		return errors.Wrap(err, "cannot comment out the squash commit message")
	}
	return script.RunCommand("git", args...)
}

// ShouldAutomaticallyAbortOnError returns whether this step should cause the command to
// automatically abort if it errors.
func (step *SquashMergeBranchStep) ShouldAutomaticallyAbortOnError() bool {
	return true
}
