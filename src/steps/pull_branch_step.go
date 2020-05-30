package steps

import "github.com/git-town/git-town/src/git"

// PullBranchStep pulls the branch with the given name from the origin remote
type PullBranchStep struct {
	NoOpStep
	BranchName string
}

// Run executes this step.
func (step *PullBranchStep) Run(repo *git.ProdRepo) error {
	return repo.Logging.Pull()
}
