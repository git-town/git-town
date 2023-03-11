package validate

import (
	"fmt"

	"github.com/fatih/color"
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
	newMainBranch, err := dialog.Select(dialog.SelectArgs{
		Options: localBranches,
		Message: mainBranchPrompt(oldMainBranch),
		Default: oldMainBranch,
	})
	if err != nil {
		return err
	}
	return repo.Config.SetMainBranch(newMainBranch)
}

func mainBranchPrompt(mainBranch string) string {
	result := "Please specify the main development branch:"
	if mainBranch != "" {
		coloredBranch := color.New(color.Bold).Add(color.FgCyan).Sprintf(mainBranch)
		result += fmt.Sprintf(" (current value: %s)", coloredBranch)
	}
	return result
}

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
