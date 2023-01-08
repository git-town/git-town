package cmd

import (
	"errors"

	"github.com/git-town/git-town/v7/src/cli"
	"github.com/git-town/git-town/v7/src/userinput"
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
		if !prodRepo.Config.IsFeatureBranch(branchName) {
			cli.Exit(errors.New("only feature branches can have parent branches"))
		}
		defaultParentBranch := prodRepo.Config.ParentBranch(branchName)
		if defaultParentBranch == "" {
			defaultParentBranch = prodRepo.Config.MainBranch()
		}
		err = prodRepo.Config.RemoveParentBranch(branchName)
		if err != nil {
			cli.Exit(err)
		}
		err = userinput.AskForBranchAncestry(branchName, defaultParentBranch, prodRepo)
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
