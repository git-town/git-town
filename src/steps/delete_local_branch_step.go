package steps

import (
	"github.com/git-town/git-town/src/git"
	"github.com/git-town/git-town/src/script"
)

// DeleteLocalBranchStep deletes the branch with the given name,
// optionally in a safe or unsafe way.
type DeleteLocalBranchStep struct {
	NoOpStep
	BranchName string
	Force      bool

	branchSha string
}

// CreateUndoStep returns the undo step for this step.
func (step *DeleteLocalBranchStep) CreateUndoStep() Step {
	return &CreateBranchStep{BranchName: step.BranchName, StartingPoint: step.branchSha}
}

// Run executes this step.
func (step *DeleteLocalBranchStep) Run() error {
	step.branchSha = git.GetBranchSha(step.BranchName)
	op := "-d"
	if step.Force || git.DoesBranchHaveUnmergedCommits(step.BranchName) {
		op = "-D"
	}
	return script.RunCommand("git", "branch", op, step.BranchName)
}
