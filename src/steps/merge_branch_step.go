package steps

import (
	"github.com/git-town/git-town/src/git"
	"github.com/git-town/git-town/src/script"
)

// MergeBranchStep merges the branch with the given name into the current branch
type MergeBranchStep struct {
	NoOpStep
	BranchName string

	previousSha string
}

// CreateAbortStep returns the abort step for this step.
func (step *MergeBranchStep) CreateAbortStep() Step {
	return &AbortMergeBranchStep{}
}

// CreateContinueStep returns the continue step for this step.
func (step *MergeBranchStep) CreateContinueStep() Step {
	return &ContinueMergeBranchStep{}
}

// CreateUndoStep returns the undo step for this step.
func (step *MergeBranchStep) CreateUndoStep() Step {
	return &ResetToShaStep{Hard: true, Sha: step.previousSha}
}

// Run executes this step.
func (step *MergeBranchStep) Run() error {
	step.previousSha = git.GetCurrentSha()
	return script.RunCommand("git", "merge", "--no-edit", step.BranchName)
}
