package dialog

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/git-town/git-town/v11/src/git"
	"github.com/git-town/git-town/v11/src/git/gitdomain"
	"github.com/git-town/git-town/v11/src/sync/syncdomain"
)

// EnterPerennialBranches lets the user update the perennial branches.
// This includes asking the user and updating the respective settings based on the user selection.
func EnterPerennialBranches(backend *git.BackendCommands, branches syncdomain.Branches) (configdomain.BranchTypes, error) {
	localBranchesWithoutMain := branches.All.LocalBranches().Remove(branches.Types.MainBranch)
	newPerennialBranchNames, err := MultiSelect(MultiSelectArgs{
		Options:  localBranchesWithoutMain.Names().Strings(),
		Defaults: branches.Types.PerennialBranches.Strings(),
		Message:  perennialBranchesPrompt(branches.Types.PerennialBranches),
	})
	if err != nil {
		return branches.Types, err
	}
	newPerennialBranches := gitdomain.NewLocalBranchNames(newPerennialBranchNames...)
	branches.Types.PerennialBranches = newPerennialBranches
	err = backend.GitTown.SetPerennialBranches(newPerennialBranches)
	return branches.Types, err
}

func perennialBranchesPrompt(perennialBranches gitdomain.LocalBranchNames) string {
	result := "Please specify perennial branches:"
	if len(perennialBranches) > 0 {
		coloredBranches := color.New(color.Bold).Add(color.FgCyan).Sprintf(perennialBranches.Join(", "))
		result += fmt.Sprintf(" (current value: %s)", coloredBranches)
	}
	return result
}
