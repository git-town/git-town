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
		ConfigureMainBranch()
		ConfigurePerennialBranches()
	}
}

// ConfigureMainBranch has the user to confgure the main branch
func ConfigureMainBranch() {
	printConfigurationHeader()
	newMainBranch := askForBranch(askForBranchOptions{
		branchNames:       git.GetLocalBranches(),
		prompt:            getMainBranchPrompt(),
		defaultBranchName: git.GetMainBranch(),
	})
	git.SetMainBranch(newMainBranch)
}

// ConfigurePerennialBranches has the user to confgure the perennial branches
func ConfigurePerennialBranches() {
	printConfigurationHeader()
	newPerennialBranches := askForBranches(askForBranchesOptions{
		branchNames:        git.GetLocalBranchesWithoutMain(),
		prompt:             getPerennialBranchesPrompt(),
		defaultBranchNames: git.GetPerennialBranches(),
	})
	git.SetPerennialBranches(newPerennialBranches)
}

// Helpers

var configurationHeaderShown bool

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

func printConfigurationHeader() {
	if !configurationHeaderShown {
		configurationHeaderShown = true
		fmt.Println("Git Town needs to be configured")
		fmt.Println()
	}
}
