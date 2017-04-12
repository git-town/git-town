package steps

import (
	"github.com/Originate/git-town/lib/git"
	"github.com/Originate/git-town/lib/script"
)

// MergeBranchStep merges the branch with the given name into the current branch
type MergeBranchStep struct {
	BranchName string
}

// CreateAbortStep returns the abort step for this step.
func (step MergeBranchStep) CreateAbortStep() Step {
	return AbortMergeBranchStep{}
}

// CreateContinueStep returns the continue step for this step.
func (step MergeBranchStep) CreateContinueStep() Step {
	return ContinueMergeBranchStep{}
}

// CreateUndoStep returns the undo step for this step.
func (step MergeBranchStep) CreateUndoStep() Step {
	return ResetToShaStep{Hard: true, Sha: git.GetCurrentSha()}
}

// Run executes this step.
func (step MergeBranchStep) Run() error {
	return script.RunCommand("git", "merge", "--no-edit", step.BranchName)
}
