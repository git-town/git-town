package steps

import (
	"github.com/git-town/git-town/src/script"
)

// DiffParentBranchStep merges the branch with the given name into the current branch
type DiffParentBranchStep struct {
	NoOpStep
	BranchName   string
	ParentBranch string
}

// Run executes this step.
func (step *DiffParentBranchStep) Run() error {
	return script.RunCommand("git", "diff", step.ParentBranch+".."+step.BranchName)
}
