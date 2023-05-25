package validate

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/git-town/git-town/v9/src/dialog"
	"github.com/git-town/git-town/v9/src/git"
)

// EnterMainBranch lets the user select a new main branch for this repo.
// This includes asking the user and updating the respective setting.
func EnterMainBranch(backend *git.BackendCommands) (string, error) {
	localBranches, err := backend.LocalBranches()
	if err != nil {
		return "", err
	}
	oldMainBranch := backend.Config.MainBranch()
	newMainBranch, err := dialog.Select(dialog.SelectArgs{
		Options: localBranches,
		Message: mainBranchPrompt(oldMainBranch),
		Default: oldMainBranch,
	})
	if err != nil {
		return "", err
	}
	return newMainBranch, backend.Config.SetMainBranch(newMainBranch)
}

func mainBranchPrompt(mainBranch string) string {
	result := "Please specify the main development branch:"
	if mainBranch != "" {
		coloredBranch := color.New(color.Bold).Add(color.FgCyan).Sprintf(mainBranch)
		result += fmt.Sprintf(" (current value: %s)", coloredBranch)
	}
	return result
}
