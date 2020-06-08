package cmd

import (
	"fmt"
	"sort"
	"strings"

	"github.com/git-town/git-town/src/cli"
	"github.com/git-town/git-town/src/git"
	"github.com/git-town/git-town/src/prompt"
	"github.com/git-town/git-town/src/util"
	"github.com/spf13/cobra"
)

var configCommand = &cobra.Command{
	Use:   "config",
	Short: "Displays your Git Town configuration",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println()
		cli.PrintLabelAndValue("Main branch", printableMainBranch(prodRepo.GetMainBranch()))
		cli.PrintLabelAndValue("Perennial branches", printablePerennialBranches(prodRepo.GetPerennialBranches()))
		mainBranch := git.Config().GetMainBranch()
		if mainBranch != "" {
			cli.PrintLabelAndValue("Branch Ancestry", printableBranchAncestry(prodRepo))
		}
		cli.PrintLabelAndValue("Pull branch strategy", git.Config().GetPullBranchStrategy())
		cli.PrintLabelAndValue("New Branch Push Flag", printableNewBranchPushFlag(prodRepo.ShouldNewBranchPush()))
	},
	Args: cobra.NoArgs,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return git.ValidateIsRepository()
	},
}

var resetConfigCommand = &cobra.Command{
	Use:   "reset",
	Short: "Resets your Git Town configuration",
	Run: func(cmd *cobra.Command, args []string) {
		git.Config().RemoveLocalGitConfiguration()
	},
	Args: cobra.NoArgs,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return git.ValidateIsRepository()
	},
}

var setupConfigCommand = &cobra.Command{
	Use:   "setup",
	Short: "Prompts to setup your Git Town configuration",
	Run: func(cmd *cobra.Command, args []string) {
		err := prompt.ConfigureMainBranch(prodRepo)
		if err != nil {
			cli.Exit(err)
		}
		err = prompt.ConfigurePerennialBranches(prodRepo)
		if err != nil {
			cli.Exit(err)
		}
	},
	Args: cobra.NoArgs,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return git.ValidateIsRepository()
	},
}

// printableBranchAncestry provides the branch ancestry in CLI printable format.
func printableBranchAncestry(repo *git.ProdRepo) string {
	roots := getBranchAncestryRoots(repo.Configuration)
	trees := make([]string, len(roots))
	for r := range roots {
		trees[r] = printableBranchTree(roots[r], repo.Configuration)
	}
	return strings.Join(trees, "\n\n")
}

// getBranchAncestryRoots returns the branches with children and no parents.
func getBranchAncestryRoots(config *git.Configuration) []string {
	parentMap := config.GetParentBranchMap()
	roots := []string{}
	for _, parent := range parentMap {
		if _, ok := parentMap[parent]; !ok && !util.DoesStringArrayContain(roots, parent) {
			roots = append(roots, parent)
		}
	}
	sort.Strings(roots)
	return roots
}

// printableBranchTree returns a user printable branch tree.
func printableBranchTree(branchName string, config *git.Configuration) (result string) {
	result += branchName
	childBranches := config.GetChildBranches(branchName)
	sort.Strings(childBranches)
	for _, childBranch := range childBranches {
		result += "\n" + util.Indent(printableBranchTree(childBranch, config))
	}
	return
}

func init() {
	configCommand.AddCommand(resetConfigCommand)
	configCommand.AddCommand(setupConfigCommand)
	RootCmd.AddCommand(configCommand)
}
