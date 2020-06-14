package prompt

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/git-town/git-town/src/git"
)

// EnsureIsConfigured has the user to confgure the main branch and perennial branches if needed.
func EnsureIsConfigured(repo *git.ProdRepo) error {
	if git.Config().GetMainBranch() == "" {
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
	newMainBranch := askForBranch(askForBranchOptions{
		branchNames:       localBranches,
		prompt:            getMainBranchPrompt(),
		defaultBranchName: git.Config().GetMainBranch(),
	})
	return git.Config().SetMainBranch(newMainBranch)
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
	newPerennialBranches := askForBranches(askForBranchesOptions{
		branchNames:        branchNames,
		prompt:             getPerennialBranchesPrompt(),
		defaultBranchNames: git.Config().GetPerennialBranches(),
	})
	return git.Config().SetPerennialBranches(newPerennialBranches)
}

// Helpers

func getMainBranchPrompt() (result string) {
	result += "Please specify the main development branch:"
	currentMainBranch := git.Config().GetMainBranch()
	if currentMainBranch != "" {
		coloredBranchName := color.New(color.Bold).Add(color.FgCyan).Sprintf(currentMainBranch)
		result += fmt.Sprintf(" (current value: %s)", coloredBranchName)
	}
	return
}

func getPerennialBranchesPrompt() (result string) {
	result += "Please specify perennial branches:"
	currentPerennialBranches := git.Config().GetPerennialBranches()
	if len(currentPerennialBranches) > 0 {
		coloredBranchNames := color.New(color.Bold).Add(color.FgCyan).Sprintf(strings.Join(currentPerennialBranches, ", "))
		result += fmt.Sprintf(" (current value: %s)", coloredBranchNames)
	}
	return
}
