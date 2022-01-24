package userinput

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/git-town/git-town/v7/src/git"
)

// EnsureIsConfigured has the user to confgure the main branch and perennial branches if needed.
func EnsureIsConfigured(repo *git.ProdRepo) error {
	if repo.Config.MainBranch() == "" {
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
		prompt:            mainBranchPrompt(repo),
		defaultBranchName: repo.Config.MainBranch(),
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
		prompt:             perennialBranchesPrompt(repo),
		defaultBranchNames: repo.Config.PerennialBranches(),
	})
	if err != nil {
		return err
	}
	return repo.Config.SetPerennialBranches(newPerennialBranches)
}

// Helpers

func mainBranchPrompt(repo *git.ProdRepo) (result string) {
	result += "Please specify the main development branch:"
	currentMainBranch := repo.Config.MainBranch()
	if currentMainBranch != "" {
		coloredBranchName := color.New(color.Bold).Add(color.FgCyan).Sprintf(currentMainBranch)
		result += fmt.Sprintf(" (current value: %s)", coloredBranchName)
	}
	return
}

func perennialBranchesPrompt(repo *git.ProdRepo) (result string) {
	result += "Please specify perennial branches:"
	currentPerennialBranches := repo.Config.PerennialBranches()
	if len(currentPerennialBranches) > 0 {
		coloredBranchNames := color.New(color.Bold).Add(color.FgCyan).Sprintf(strings.Join(currentPerennialBranches, ", "))
		result += fmt.Sprintf(" (current value: %s)", coloredBranchNames)
	}
	return
}
