package cmd

import (
	"github.com/git-town/git-town/src/git"
	"github.com/git-town/git-town/src/prompt"
	"github.com/git-town/git-town/src/util"
	"github.com/spf13/cobra"
)

var setParentBranchCommand = &cobra.Command{
	Use:   "set-parent-branch",
	Short: "Prompts to set the parent branch for the current branch",
	Long:  `Prompts to set the parent branch for the current branch`,
	Run: func(cmd *cobra.Command, args []string) {
		branchName := git.GetCurrentBranchName()
		git.Config().EnsureIsFeatureBranch(branchName, "Only feature branches can have parent branches.")
		defaultParentBranch := git.Config().GetParentBranch(branchName)
		if defaultParentBranch == "" {
			defaultParentBranch = git.Config().GetMainBranch()
		}
		git.Config().DeleteParentBranch(branchName)
		prompt.AskForBranchAncestry(branchName, defaultParentBranch)
	},
	Args: cobra.NoArgs,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return util.FirstError(
			git.ValidateIsRepository,
			validateIsConfigured,
		)
	},
}

func init() {
	RootCmd.AddCommand(setParentBranchCommand)
}
