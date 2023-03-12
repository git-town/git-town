package validate

import (
	"fmt"

	"github.com/git-town/git-town/v7/src/git"
)

// IsConfigured is a validationCondition that verifies that the given Git repo contains necessary Git Town configuration.
func IsConfigured(repo *git.ProdRepo) error {
	if repo.Config.MainBranch() == "" {
		fmt.Print("Git Town needs to be configured\n\n")
		err := EnterMainBranch(repo)
		if err != nil {
			return err
		}
		return EnterPerennialBranches(repo)
	}
	return repo.RemoveOutdatedConfiguration()
}
