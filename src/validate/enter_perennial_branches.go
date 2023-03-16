package validate

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/git-town/git-town/v7/src/dialog"
	"github.com/git-town/git-town/v7/src/git"
)

// EnterPerennialBranches lets the user update the perennial branches.
// This includes asking the user and updating the respective settings based on the user selection.
func EnterPerennialBranches(repo *git.PublicRepo, mainBranch string) error {
	localBranchesWithoutMain, err := repo.LocalBranchesWithoutMain(mainBranch)
	if err != nil {
		return err
	}
	oldPerennialBranches := repo.Config.PerennialBranches()
	newPerennialBranches, err := dialog.MultiSelect(dialog.MultiSelectArgs{
		Options:  localBranchesWithoutMain,
		Defaults: oldPerennialBranches,
		Message:  perennialBranchesPrompt(oldPerennialBranches),
	})
	if err != nil {
		return err
	}
	return repo.Config.SetPerennialBranches(newPerennialBranches)
}

func perennialBranchesPrompt(perennialBranches []string) string {
	result := "Please specify perennial branches:"
	if len(perennialBranches) > 0 {
		coloredBranches := color.New(color.Bold).Add(color.FgCyan).Sprintf(strings.Join(perennialBranches, ", "))
		result += fmt.Sprintf(" (current value: %s)", coloredBranches)
	}
	return result
}
