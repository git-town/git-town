package validate

import (
	"github.com/git-town/git-town/v7/src/dialog"
	"github.com/git-town/git-town/v7/src/git"
)

func ConfigurePerennialBranches(repo *git.ProdRepo) error {
	localBranchesWithoutMain, err := repo.Silent.LocalBranchesWithoutMain()
	if err != nil {
		return err
	}
	perennialBranches := repo.Config.PerennialBranches()
	newPerennialBranches, err := dialog.EnterPerennialBranches(localBranchesWithoutMain, perennialBranches)
	if err != nil {
		return err
	}
	return repo.Config.SetPerennialBranches(newPerennialBranches)
}
