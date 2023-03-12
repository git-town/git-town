package validate

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/git-town/git-town/v7/src/dialog"
	"github.com/git-town/git-town/v7/src/git"
)

func EnterPerennialBranches(repo *git.ProdRepo) error {
	localBranchesWithoutMain, err := repo.Silent.LocalBranchesWithoutMain()
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
