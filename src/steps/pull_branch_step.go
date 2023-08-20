package steps

import (
	"github.com/git-town/git-town/v9/src/git"
	"github.com/git-town/git-town/v9/src/hosting"
)

// PullBranchStep updates the branch with the given name with commits from its remote.
// TODO: rename to PullCurrentBranchStep and remove the "Branch" field.
type PullBranchStep struct {
	Branch string
	EmptyStep
}

func (step *PullBranchStep) Run(run *git.ProdRunner, _ hosting.Connector) error {
	return run.Frontend.Pull()
}
