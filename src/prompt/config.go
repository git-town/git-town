package prompt

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/git-town/git-town/src/git"
)

// EnsureIsConfigured has the user to confgure the main branch and perennial branches if needed.
func EnsureIsConfigured(repo *git.ProdRepo) error {
	if repo.Config.GetMainBranch() == "" {
		fmt.Println("Git Town needs to be configured")
		fmt.Println()
		err := ConfigureMainBranch(repo)
		if err != nil {
			return err
		}
		return ConfigurePerennialBranches(repo)
	}
	return nil
}

// ConfigureMainBranch has the user to confgure the main branch.
func ConfigureMainBranch(repo *git.ProdRepo) error {
	localBranches, err := repo.Silent.LocalBranches()
	if err != nil {
		return err
	}
	newMainBranch, err := askForBranch(askForBranchOptions{
		branchNames:       localBranches,
		prompt:            getMainBranchPrompt(repo),
		defaultBranchName: repo.Config.GetMainBranch(),
	})
	if err != nil {
		return err
	}
	return repo.Config.SetMainBranch(newMainBranch)
}

// ConfigurePerennialBranches has the user to confgure the perennial branches.
func ConfigurePerennialBranches(repo *git.ProdRepo) error {
	branchNames, err := repo.Silent.LocalBranchesWithoutMain()
	if err != nil {
		return err
	}
	if len(branchNames) == 0 {
		return nil
	}
	newPerennialBranches, err := askForBranches(askForBranchesOptions{
		branchNames:        branchNames,
		prompt:             getPerennialBranchesPrompt(repo),
		defaultBranchNames: repo.Config.GetPerennialBranches(),
	})
	if err != nil {
		return err
	}
	return repo.Config.SetPerennialBranches(newPerennialBranches)
}

// Helpers

func getMainBranchPrompt(repo *git.ProdRepo) (result string) {
	result += "Please specify the main development branch:"
	currentMainBranch := repo.Config.GetMainBranch()
	if currentMainBranch != "" {
		coloredBranchName := color.New(color.Bold).Add(color.FgCyan).Sprintf(currentMainBranch)
		result += fmt.Sprintf(" (current value: %s)", coloredBranchName)
	}
	return
}

func getPerennialBranchesPrompt(repo *git.ProdRepo) (result string) {
	result += "Please specify perennial branches:"
	currentPerennialBranches := repo.Config.GetPerennialBranches()
	if len(currentPerennialBranches) > 0 {
		coloredBranchNames := color.New(color.Bold).Add(color.FgCyan).Sprintf(strings.Join(currentPerennialBranches, ", "))
		result += fmt.Sprintf(" (current value: %s)", coloredBranchNames)
	}
	return
}
