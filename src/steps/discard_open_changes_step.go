package steps

import (
	"github.com/git-town/git-town/src/git"
)

// DiscardOpenChangesStep resets the branch to the last commit, discarding uncommitted changes.
type DiscardOpenChangesStep struct {
	NoOpStep
}

// Run executes this step.
func (step *DiscardOpenChangesStep) Run(repo *git.ProdRepo) error {
	return repo.Logging.DiscardOpenChanges()
}
