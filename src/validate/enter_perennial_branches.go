package validate

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/git-town/git-town/v9/src/config"
	"github.com/git-town/git-town/v9/src/dialog"
	"github.com/git-town/git-town/v9/src/git"
)

// EnterPerennialBranches lets the user update the perennial branches.
// This includes asking the user and updating the respective settings based on the user selection.
func EnterPerennialBranches(backend *git.BackendCommands, allBranches git.BranchesSyncStatus, oldBranches config.BranchDurations) (config.BranchDurations, error) {
	localBranchesWithoutMain := allBranches.LocalBranches().Remove(oldBranches.MainBranch).BranchNames()
	newPerennialBranches, err := dialog.MultiSelect(dialog.MultiSelectArgs{
		Options:  localBranchesWithoutMain,
		Defaults: oldBranches.PerennialBranches,
		Message:  perennialBranchesPrompt(oldBranches.PerennialBranches),
	})
	if err != nil {
		return oldBranches, err
	}
	oldBranches.PerennialBranches = newPerennialBranches
	err = backend.Config.SetPerennialBranches(newPerennialBranches)
	return oldBranches, err
}

func perennialBranchesPrompt(perennialBranches []string) string {
	result := "Please specify perennial branches:"
	if len(perennialBranches) > 0 {
		coloredBranches := color.New(color.Bold).Add(color.FgCyan).Sprintf(strings.Join(perennialBranches, ", "))
		result += fmt.Sprintf(" (current value: %s)", coloredBranches)
	}
	return result
}
