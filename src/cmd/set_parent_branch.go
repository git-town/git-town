package cmd

import (
	"fmt"
	"os"

	"github.com/git-town/git-town/src/cli"
	"github.com/git-town/git-town/src/git"
	"github.com/git-town/git-town/src/prompt"
	"github.com/spf13/cobra"
)

var setParentBranchCommand = &cobra.Command{
	Use:   "set-parent-branch",
	Short: "Prompts to set the parent branch for the current branch",
	Long:  `Prompts to set the parent branch for the current branch`,
	Run: func(cmd *cobra.Command, args []string) {
		branchName, err := prodRepo.Silent.CurrentBranch()
		if err != nil {
			cli.Exit(err)
		}
		if !git.Config().IsFeatureBranch(branchName) {
			fmt.Println("Error: only feature branches can have parent branches")
			os.Exit(1)
		}
		defaultParentBranch := git.Config().GetParentBranch(branchName)
		if defaultParentBranch == "" {
			defaultParentBranch = git.Config().GetMainBranch()
		}
		err = git.Config().DeleteParentBranch(branchName)
		if err != nil {
			cli.Exit(err)
		}
		err = prompt.AskForBranchAncestry(branchName, defaultParentBranch, prodRepo)
		if err != nil {
			cli.Exit(err)
		}
	},
	Args: cobra.NoArgs,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := ValidateIsRepository(prodRepo); err != nil {
			return err
		}
		return validateIsConfigured(prodRepo)
	},
}

func init() {
	RootCmd.AddCommand(setParentBranchCommand)
}
