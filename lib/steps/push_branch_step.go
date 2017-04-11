package steps

import (
	"github.com/Originate/git-town/lib/git"
	"github.com/Originate/git-town/lib/script"
)

type PushBranchStep struct {
	NoAutomaticAbortOnError
	NoUndoStepAfterRun
	BranchName string
	Force      bool
	Undoable   bool
}

func (step PushBranchStep) CreateAbortStep() Step {
	return NoOpStep{}
}

func (step PushBranchStep) CreateContinueStep() Step {
	return NoOpStep{}
}

func (step PushBranchStep) CreateUndoStepBeforeRun() Step {
	if step.Undoable {
		return PushBranchAfterCurrentBranchSteps{}
	} else {
		return SkipCurrentBranchSteps{}
	}
}

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
