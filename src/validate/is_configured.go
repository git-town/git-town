package validate

import (
	"fmt"

	"github.com/git-town/git-town/v7/src/git"
)

// IsConfigured verifies that the given Git repo contains necessary Git Town configuration.
func IsConfigured(repo *git.PublicRepo) error {
	mainBranch := repo.Config.MainBranch()
	if mainBranch == "" {
		fmt.Print("Git Town needs to be configured\n\n")
		err := EnterMainBranch(repo)
		if err != nil {
			return err
		}
		return EnterPerennialBranches(repo, mainBranch)
	}
	return repo.RemoveOutdatedConfiguration()
}
