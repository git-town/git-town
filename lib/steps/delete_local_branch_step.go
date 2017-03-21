package steps

import (
	"github.com/Originate/git-town/lib/git"
	"github.com/Originate/git-town/lib/script"
)

type DeleteLocalBranchStep struct {
	BranchName string
	Force      bool
}

func (step DeleteLocalBranchStep) CreateAbortStep() Step {
	return NoOpStep{}
}

func (step DeleteLocalBranchStep) CreateContinueStep() Step {
	return NoOpStep{}
}

func (step DeleteLocalBranchStep) CreateUndoStep() Step {
	return NoOpStep{}
}

func (step DeleteLocalBranchStep) Run() error {
	op := "-d"
	if step.Force || git.DoesBranchHaveUnmergedCommits(step.BranchName) {
		op = "-D"
	}
	return script.RunCommand("git", "branch", op, step.BranchName)
}
