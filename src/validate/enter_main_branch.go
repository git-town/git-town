package validate

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/git-town/git-town/v7/src/dialog"
	"github.com/git-town/git-town/v7/src/git"
)

// EnterMainBranch lets the user select a new main branch for this repo.
// This includes asking the user and updating the respective setting.
func EnterMainBranch(repo *git.PublicRepo) (string, error) {
	localBranches, err := repo.LocalBranches()
	if err != nil {
		return "", err
	}
	oldMainBranch := repo.Config.MainBranch()
	newMainBranch, err := dialog.Select(dialog.SelectArgs{
		Options: localBranches,
		Message: mainBranchPrompt(oldMainBranch),
		Default: oldMainBranch,
	})
	if err != nil {
		return "", err
	}
	return newMainBranch, repo.Config.SetMainBranch(newMainBranch)
}

func mainBranchPrompt(mainBranch string) string {
	result := "Please specify the main development branch:"
	if mainBranch != "" {
		coloredBranch := color.New(color.Bold).Add(color.FgCyan).Sprintf(mainBranch)
		result += fmt.Sprintf(" (current value: %s)", coloredBranch)
	}
	return result
}
