package steps

import (
	"github.com/git-town/git-town/src/git"
	"github.com/git-town/git-town/src/script"
)

// CreateRemoteBranchStep pushes the current branch up to origin.
type CreateRemoteBranchStep struct {
	NoOpStep
	BranchName string
	Sha        string
}

// Run executes this step.
func (step *CreateRemoteBranchStep) Run(repo *git.ProdRepo) error {
	return script.RunCommand("git", "push", "origin", step.Sha+":refs/heads/"+step.BranchName)
}
