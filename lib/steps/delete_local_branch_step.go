package steps

import (
	"github.com/Originate/git-town/lib/git"
	"github.com/Originate/git-town/lib/script"
)

// DeleteLocalBranchStep deletes the branch with the given name,
// optionally in a safe or unsafe way.
type DeleteLocalBranchStep struct {
	NoAutomaticAbortOnError
	NoUndoStepAfterRun
	BranchName string
	Force      bool
}

// CreateAbortStep returns the abort step for this step.
func (step DeleteLocalBranchStep) CreateAbortStep() Step {
	return NoOpStep{}
}

// CreateContinueStep returns the continue step for this step.
func (step DeleteLocalBranchStep) CreateContinueStep() Step {
	return NoOpStep{}
}

// CreateUndoStepBeforeRun returns the undo step for this step before it is run.
func (step DeleteLocalBranchStep) CreateUndoStepBeforeRun() Step {
	sha := git.GetBranchSha(step.BranchName)
	return CreateBranchStep{BranchName: step.BranchName, StartingPoint: sha}
}

// Run executes this step.
func (step DeleteLocalBranchStep) Run() error {
	op := "-d"
	if step.Force || git.DoesBranchHaveUnmergedCommits(step.BranchName) {
		op = "-D"
	}
	return script.RunCommand("git", "branch", op, step.BranchName)
}
