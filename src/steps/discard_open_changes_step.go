package steps

import (
	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/hosting"
)

// DiscardOpenChangesStep resets the branch to the last commit, discarding uncommitted changes.
type DiscardOpenChangesStep struct {
	NoOpStep
}

func (step *DiscardOpenChangesStep) Run(repo *git.ProdRepo, driver hosting.Driver) error {
	return repo.Logging.DiscardOpenChanges()
}
