package dialog

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/git"
)

// EnterPerennialBranches lets the user update the perennial branches.
// This includes asking the user and updating the respective settings based on the user selection.
func EnterPerennialBranches(backend *git.BackendCommands, allBranches domain.BranchInfos, branchTypes domain.BranchTypes) (domain.BranchTypes, error) {
	localBranchesWithoutMain := allBranches.LocalBranches().Remove(branchTypes.MainBranch)
	newPerennialBranchNames, err := MultiSelect(MultiSelectArgs{
		Options:  localBranchesWithoutMain.Names().Strings(),
		Defaults: branchTypes.PerennialBranches.Strings(),
		Message:  perennialBranchesPrompt(branchTypes.PerennialBranches),
	})
	if err != nil {
		return branchTypes, err
	}
	newPerennialBranches := domain.NewLocalBranchNames(newPerennialBranchNames...)
	branchTypes.PerennialBranches = newPerennialBranches
	err = backend.Config.SetPerennialBranches(newPerennialBranches)
	return branchTypes, err
}

func perennialBranchesPrompt(perennialBranches domain.LocalBranchNames) string {
	result := "Please specify perennial branches:"
	if len(perennialBranches) > 0 {
		coloredBranches := color.New(color.Bold).Add(color.FgCyan).Sprintf(perennialBranches.Join(", "))
		result += fmt.Sprintf(" (current value: %s)", coloredBranches)
	}
	return result
}
