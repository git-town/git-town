package prompt

import (
	"fmt"
	"strings"

	"github.com/Originate/git-town/src/git"
	"github.com/fatih/color"
)

// EnsureIsConfigured has the user to confgure the main branch and perennial branches if needed
func EnsureIsConfigured() {
	if git.GetMainBranch() == "" {
		fmt.Println("Git Town needs to be configured")
		fmt.Println()
		ConfigureMainBranch()
		ConfigurePerennialBranches()
	}
}

// ConfigureMainBranch has the user to confgure the main branch
func ConfigureMainBranch() {
	newMainBranch := askForBranch(askForBranchOptions{
		branchNames:       git.GetLocalBranches(),
		prompt:            getMainBranchPrompt(),
		defaultBranchName: git.GetMainBranch(),
	})
	git.SetMainBranch(newMainBranch)
}

// ConfigurePerennialBranches has the user to confgure the perennial branches
func ConfigurePerennialBranches() {
	branchNames := git.GetLocalBranchesWithoutMain()
	if len(branchNames) == 0 {
		return
	}
	newPerennialBranches := askForBranches(askForBranchesOptions{
		branchNames:        branchNames,
		prompt:             getPerennialBranchesPrompt(),
		defaultBranchNames: git.GetPerennialBranches(),
	})
	git.SetPerennialBranches(newPerennialBranches)
}

// Helpers

func getMainBranchPrompt() (result string) {
	result += "Please specify the main development branch:"
	currentMainBranch := git.GetMainBranch()
	if currentMainBranch != "" {
		coloredBranchName := color.New(color.Bold).Add(color.FgCyan).Sprintf(currentMainBranch)
		result += fmt.Sprintf(" (current value: %s)", coloredBranchName)
	}
	return
}

func getPerennialBranchesPrompt() (result string) {
	result += "Please specify perennial branches:"
	currentPerennialBranches := git.GetPerennialBranches()
	if len(currentPerennialBranches) > 0 {
		coloredBranchNames := color.New(color.Bold).Add(color.FgCyan).Sprintf(strings.Join(currentPerennialBranches, ", "))
		result += fmt.Sprintf(" (current value: %s)", coloredBranchNames)
	}
	return
}
