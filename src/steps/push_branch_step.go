package steps

import (
	"github.com/git-town/git-town/src/dryrun"
	"github.com/git-town/git-town/src/git"
	"github.com/git-town/git-town/src/script"
)

// PushBranchStep pushes the branch with the given name to the origin remote.
// Optionally with force.
type PushBranchStep struct {
	NoOpStep
	BranchName string
	Force      bool
	Undoable   bool
}

// CreateUndoStep returns the undo step for this step.
func (step *PushBranchStep) CreateUndoStep() Step {
	if step.Undoable {
		return &PushBranchAfterCurrentBranchSteps{}
	}
	return &SkipCurrentBranchSteps{}
}

// Run executes this step.
func (step *PushBranchStep) Run() error {
	if !git.ShouldBranchBePushed(step.BranchName) && !dryrun.IsActive() {
		return nil
	}
	if step.Force {
		return script.RunCommand("git", "push", "-f", "origin", step.BranchName)
	}
	if git.GetCurrentBranchName() == step.BranchName {
		return script.RunCommand("git", "push")
	}
	return script.RunCommand("git", "push", "origin", step.BranchName)
}
