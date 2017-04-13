package steps

import (
	"log"

	"github.com/Originate/git-town/lib/git"
	"github.com/Originate/git-town/lib/prompt"
	"github.com/Originate/git-town/lib/script"
	"github.com/Originate/git-town/lib/util"
)

// SquashMergeBranchStep squash merges the branch with the given name into the current branch
type SquashMergeBranchStep struct {
	NoContinueStep
	NoUndoStepBeforeRun
	BranchName    string
	CommitMessage string
}

// CreateAbortStep returns the abort step for this step.
func (step SquashMergeBranchStep) CreateAbortStep() Step {
	return DiscardOpenChangesStep{}
}

// CreateUndoStepAfterRun returns the undo step for this step after it is run.
func (step SquashMergeBranchStep) CreateUndoStepAfterRun() Step {
	return RevertCommitStep{Sha: git.GetCurrentSha()}
}

// GetAutomaticAbortErrorMessage returns the error message to display when this step
// cause the command to automatically abort.
func (step SquashMergeBranchStep) GetAutomaticAbortErrorMessage() string {
	return "Aborted because commit exited with error"
}

// Run executes this step.
func (step SquashMergeBranchStep) Run() error {
	err := script.RunCommand("git", "merge", "--squash", step.BranchName)
	if err != nil {
		log.Fatal("Error squash merging:", err)
	}
	commitCmd := []string{"git", "commit"}
	if step.CommitMessage != "" {
		commitCmd = append(commitCmd, "-m", step.CommitMessage)
	}
	author := prompt.GetSquashCommitAuthor(step.BranchName)
	if author != git.GetLocalAuthor() {
		commitCmd = append(commitCmd, "--author", author)
	}
	util.GetCommandOutput("sed", "-i", "-e", "s/^/# /g", ".git/SQUASH_MSG")
	return script.RunCommand(commitCmd...)
}

// ShouldAutomaticallyAbortOnError returns whether this step should cause the command to
// automatically abort if it errors.
func (step SquashMergeBranchStep) ShouldAutomaticallyAbortOnError() bool {
	return true
}
