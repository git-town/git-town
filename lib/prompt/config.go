package prompt

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Originate/git-town/lib/git"
	"github.com/fatih/color"
)

// PromptForMainBranch asks the user to confgure the main branch
func PromptForMainBranch() {
	printConfigurationHeader()
	newMainBranch := askForBranch(branchPromptConfig{
		branchNames: git.GetLocalBranches(),
		prompt:      getMainBranchPrompt(),
		validate: func(branchName string) error {
			if branchName == "" {
				return errors.New("A main development branch is required to enable the features provided by Git Town")
			}
			return nil
		},
	})
	git.SetMainBranch(newMainBranch)
}

// PromptForMainBranch asks the user to confgure the perennial branches
func PromptForPerennialBranches() {
	printConfigurationHeader()
	var newPerennialBranches []string
	for {
		newPerennialBranch := askForBranch(branchPromptConfig{
			branchNames: git.GetLocalBranches(),
			prompt:      getPerennialBranchesPrompt(),
			validate: func(branchName string) error {
				if branchName == git.GetMainBranch() {
					return fmt.Errorf("'%s' is already set as the main branch", branchName)
				}
				return nil
			},
		})
		if newPerennialBranch == "" {
			break
		}
		newPerennialBranches = append(newPerennialBranches, newPerennialBranch)
	}
	git.SetPerennialBranches(newPerennialBranches)
}

// Helpers

var configurationHeaderShown bool

func getMainBranchPrompt() (result string) {
	result += "Please specify the main development branch by name or number"
	currentMainBranch := git.GetMainBranch()
	if currentMainBranch != "" {
		coloredBranchName := color.New(color.Bold).Add(color.FgCyan).Sprintf(currentMainBranch)
		result += fmt.Sprintf(" (current value: %s)", coloredBranchName)
	}
	result += ":"
	return
}

func getPerennialBranchesPrompt() (result string) {
	result += "Please specify a perennial branch by name or number. Leave it blank to finish"
	currentPerennialBranches := git.GetPerennialBranches()
	if len(currentPerennialBranches) > 0 {
		coloredBranchNames := color.New(color.Bold).Add(color.FgCyan).Sprintf(strings.Join(currentPerennialBranches, ", "))
		result += fmt.Sprintf(" (current value: %s)", coloredBranchNames)
	}
	result += ":"
	return
}

func printConfigurationHeader() {
	if !configurationHeaderShown {
		configurationHeaderShown = true
		fmt.Println("Git Town needs to be configured")
		fmt.Println()
		printNumberedBranches(git.GetLocalBranches())
		fmt.Println()
	}
}
