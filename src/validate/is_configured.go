package validate

import (
	"fmt"

	"github.com/git-town/git-town/v7/src/dialog"
	"github.com/git-town/git-town/v7/src/git"
)

// IsConfigured is a validationCondition that verifies that the given Git repo contains necessary Git Town configuration.
func IsConfigured(repo *git.ProdRepo) error {
	if repo.Config.MainBranch() == "" {
		fmt.Print("Git Town needs to be configured\n\n")
		err := ConfigureMainBranch(repo)
		if err != nil {
			return err
		}
		return ConfigurePerennialBranches(repo)
	}
	return repo.RemoveOutdatedConfiguration()
}

func ConfigureMainBranch(repo *git.ProdRepo) error {
	localBranches, err := repo.Silent.LocalBranches()
	if err != nil {
		return err
	}
	oldMainBranch := repo.Config.MainBranch()
	newMainBranch, err := dialog.AskMainBranch(oldMainBranch, localBranches)
	if err != nil {
		return err
	}
	return repo.Config.SetMainBranch(newMainBranch)
}

func ConfigurePerennialBranches(repo *git.ProdRepo) error {
	localBranchesWithoutMain, err := repo.Silent.LocalBranchesWithoutMain()
	if err != nil {
		return err
	}
	perennialBranches := repo.Config.PerennialBranches()
	newPerennialBranches, err := dialog.AskPerennialBranches(localBranchesWithoutMain, perennialBranches)
	if err != nil {
		return err
	}
	return repo.Config.SetPerennialBranches(newPerennialBranches)
}
