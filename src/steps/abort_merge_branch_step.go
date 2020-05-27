package steps

import (
	"github.com/git-town/git-town/src/git"
	"github.com/git-town/git-town/src/script"
)

// AbortMergeBranchStep aborts the current merge conflict.
type AbortMergeBranchStep struct {
	NoOpStep
}

// Run executes this step.
func (step *AbortMergeBranchStep) Run(repo *git.ProdRepo) error {
	return script.RunCommand("git", "merge", "--abort")
}
