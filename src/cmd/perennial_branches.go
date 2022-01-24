package cmd

import (
	"github.com/git-town/git-town/v7/src/cli"
	"github.com/git-town/git-town/v7/src/userinput"
	"github.com/spf13/cobra"
)

var perennialBranchesCommand = &cobra.Command{
	Use:   "perennial-branches",
	Short: "Displays your perennial branches",
	Long: `Displays your perennial branches

Perennial branches are long-lived branches.
They cannot be shipped.`,
	Run: func(cmd *cobra.Command, args []string) {
		cli.Println(cli.PrintablePerennialBranches(prodRepo.Config.PerennialBranches()))
	},
	Args: cobra.NoArgs,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return ValidateIsRepository(prodRepo)
	},
}

var updatePrennialBranchesCommand = &cobra.Command{
	Use:   "update",
	Short: "Prompts to update your perennial branches",
	Long:  `Prompts to update your perennial branches`,
	Run: func(cmd *cobra.Command, args []string) {
		err := userinput.ConfigurePerennialBranches(prodRepo)
		if err != nil {
			cli.Exit(err)
		}
	},
	Args: cobra.NoArgs,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return ValidateIsRepository(prodRepo)
	},
}

func init() {
	perennialBranchesCommand.AddCommand(updatePrennialBranchesCommand)
	RootCmd.AddCommand(perennialBranchesCommand)
}
