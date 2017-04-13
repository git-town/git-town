package steps

import (
	"github.com/Originate/git-town/lib/git"
	"github.com/Originate/git-town/lib/script"
)

// PushBranchStep pushes the branch with the given name to the origin remote.
// Optionally with force.
type PushBranchStep struct {
	NoOpStep
	BranchName string
	Force      bool
	Undoable   bool
}

// CreateUndoStepBeforeRun returns the undo step for this step before it is run.
func (step PushBranchStep) CreateUndoStepBeforeRun() Step {
	if step.Undoable {
		return PushBranchAfterCurrentBranchSteps{}
	}
	return SkipCurrentBranchSteps{}
}

// Run executes this step.
func (step PushBranchStep) Run() error {
	if !git.ShouldBranchBePushed(step.BranchName) {
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
