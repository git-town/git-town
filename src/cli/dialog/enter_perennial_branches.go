package dialog

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/git-town/git-town/v11/src/git"
	"github.com/git-town/git-town/v11/src/git/gitdomain"
)

// EnterPerennialBranches lets the user update the perennial branches.
// This includes asking the user and updating the respective settings based on the user selection.
func EnterPerennialBranches(backend *git.BackendCommands, config *configdomain.FullConfig, allBranches gitdomain.BranchInfos) error {
	localBranchesWithoutMain := allBranches.LocalBranches().Remove(config.MainBranch)
	newPerennialBranchNames, err := MultiSelect(MultiSelectArgs{
		Options:  localBranchesWithoutMain.Names().Strings(),
		Defaults: config.PerennialBranches.Strings(),
		Message:  perennialBranchesPrompt(config.PerennialBranches),
	})
	if err != nil {
		return err
	}
	newPerennialBranches := gitdomain.NewLocalBranchNames(newPerennialBranchNames...)
	config.PerennialBranches = newPerennialBranches
	return backend.Config.SetPerennialBranches(newPerennialBranches)
}

func perennialBranchesPrompt(perennialBranches gitdomain.LocalBranchNames) string {
	result := "Please specify perennial branches:"
	if len(perennialBranches) > 0 {
		coloredBranches := color.New(color.Bold).Add(color.FgCyan).Sprintf(perennialBranches.Join(", "))
		result += fmt.Sprintf(" (current value: %s)", coloredBranches)
	}
	return result
}
