package validate

import (
	"github.com/git-town/git-town/v7/src/dialog"
	"github.com/git-town/git-town/v7/src/git"
)

// IsConfigured is a validationCondition that verifies that the given Git repo contains necessary Git Town configuration.
func IsConfigured(repo *git.ProdRepo) error {
	err := dialog.EnsureIsConfigured(repo)
	if err != nil {
		return err
	}
	return repo.RemoveOutdatedConfiguration()
}
