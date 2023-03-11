package dialog

import (
	"fmt"

	"github.com/fatih/color"
)

// EnterMainBranch lets the user enter the main branch.
func EnterMainBranch(mainBranch string, localBranches []string) (string, error) {
	newMainBranch, err := EnterBranch(AskForBranchOptions{
		Branches:      localBranches,
		Prompt:        mainBranchPrompt(mainBranch),
		DefaultBranch: mainBranch,
	})
	if err != nil {
		return "", err
	}
	return newMainBranch, nil
}

func mainBranchPrompt(mainBranch string) string {
	result := "Please specify the main development branch:"
	currentMainBranch := mainBranch
	if currentMainBranch != "" {
		coloredBranch := color.New(color.Bold).Add(color.FgCyan).Sprintf(currentMainBranch)
		result += fmt.Sprintf(" (current value: %s)", coloredBranch)
	}
	return result
}
