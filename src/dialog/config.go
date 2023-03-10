package dialog

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
)

// AskMainBranch lets the user enter the main branch.
func AskMainBranch(mainBranch string, localBranches []string) (string, error) {
	newMainBranch, err := askForBranch(askForBranchOptions{
		branches:      localBranches,
		prompt:        mainBranchPrompt(mainBranch),
		defaultBranch: mainBranch,
	})
	if err != nil {
		return "", err
	}
	return newMainBranch, nil
}

// AskPerennialBranches lets the user enter the perennial branches.
func AskPerennialBranches(localBranchesWithoutMain []string, perennialBranches []string) ([]string, error) {
	if len(localBranchesWithoutMain) == 0 {
		return []string{}, nil
	}
	newPerennialBranches, err := askForBranches(askForBranchesOptions{
		branches:        localBranchesWithoutMain,
		prompt:          perennialBranchesPrompt(perennialBranches),
		defaultBranches: perennialBranches,
	})
	if err != nil {
		return []string{}, err
	}
	return newPerennialBranches, nil
}

// Helpers

func mainBranchPrompt(mainBranch string) string {
	result := "Please specify the main development branch:"
	currentMainBranch := mainBranch
	if currentMainBranch != "" {
		coloredBranch := color.New(color.Bold).Add(color.FgCyan).Sprintf(currentMainBranch)
		result += fmt.Sprintf(" (current value: %s)", coloredBranch)
	}
	return result
}

func perennialBranchesPrompt(perennialBranches []string) string {
	result := "Please specify perennial branches:"
	if len(perennialBranches) > 0 {
		coloredBranches := color.New(color.Bold).Add(color.FgCyan).Sprintf(strings.Join(perennialBranches, ", "))
		result += fmt.Sprintf(" (current value: %s)", coloredBranches)
	}
	return result
}
