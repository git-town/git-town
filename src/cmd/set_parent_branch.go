package cmd

import (
	"github.com/Originate/git-town/src/git"
	"github.com/Originate/git-town/src/prompt"
	"github.com/Originate/git-town/src/util"
	"github.com/spf13/cobra"
)

var setParentBranchCommand = &cobra.Command{
	Use:   "set-parent-branch",
	Short: "Prompts to set the parent branch for the current branch",
	Long:  `Prompts to set the parent branch for the current branch`,
	Run: func(cmd *cobra.Command, args []string) {
		promptForParentBranch()
	},
	Args: cobra.NoArgs,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return util.FirstError(
			git.ValidateIsRepository,
			validateIsConfigured,
		)
	},
}

func promptForParentBranch() {
	branchName := git.GetCurrentBranchName()
	git.EnsureIsFeatureBranch(branchName, "Only feature branches can have parent branches.")
	defaultParentBranch := git.GetParentBranch(branchName)
	if defaultParentBranch == "" {
		defaultParentBranch = git.GetMainBranch()
	}
	git.DeleteParentBranch(branchName)
	prompt.AskForBranchAncestry(branchName, defaultParentBranch)
}

func init() {
	RootCmd.AddCommand(setParentBranchCommand)
}
