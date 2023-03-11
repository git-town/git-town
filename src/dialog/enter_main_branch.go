package dialog

import (
	"fmt"

	"github.com/fatih/color"
)

// EnterMainBranch lets the user enter the main branch.
func EnterMainBranch(mainBranch string, localBranches []string) (string, error) {
	return Select(SelectArgs{
		Options: localBranches,
		Message: mainBranchPrompt(mainBranch),
		Default: mainBranch,
	})
}

func mainBranchPrompt(mainBranch string) string {
	result := "Please specify the main development branch:"
	if mainBranch != "" {
		coloredBranch := color.New(color.Bold).Add(color.FgCyan).Sprintf(mainBranch)
		result += fmt.Sprintf(" (current value: %s)", coloredBranch)
	}
	return result
}
