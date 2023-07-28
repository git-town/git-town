package dialog

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/git-town/git-town/v9/src/config"
	"github.com/git-town/git-town/v9/src/git"
)

// EnterPerennialBranches lets the user update the perennial branches.
// This includes asking the user and updating the respective settings based on the user selection.
func EnterPerennialBranches(backend *git.BackendCommands, allBranches git.BranchesSyncStatus, branches config.BranchDurations) (config.BranchDurations, error) {
	localBranchesWithoutMain := allBranches.LocalBranches().Remove(branches.MainBranch).BranchNames()
	newPerennialBranches, err := MultiSelect(MultiSelectArgs{
		Options:  localBranchesWithoutMain,
		Defaults: branches.PerennialBranches,
		Message:  perennialBranchesPrompt(branches.PerennialBranches),
	})
	if err != nil {
		return branches, err
	}
	branches.PerennialBranches = newPerennialBranches
	err = backend.Config.SetPerennialBranches(newPerennialBranches)
	return branches, err
}

func perennialBranchesPrompt(perennialBranches []string) string {
	result := "Please specify perennial branches:"
	if len(perennialBranches) > 0 {
		coloredBranches := color.New(color.Bold).Add(color.FgCyan).Sprintf(strings.Join(perennialBranches, ", "))
		result += fmt.Sprintf(" (current value: %s)", coloredBranches)
	}
	return result
}
