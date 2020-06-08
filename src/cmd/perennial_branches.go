package cmd

import (
	"strings"

	"github.com/git-town/git-town/src/cli"
	"github.com/git-town/git-town/src/git"
	"github.com/git-town/git-town/src/prompt"
	"github.com/spf13/cobra"
)

var perennialBranchesCommand = &cobra.Command{
	Use:   "perennial-branches",
	Short: "Displays your perennial branches",
	Long: `Displays your perennial branches

Perennial branches are long-lived branches.
They cannot be shipped.`,
	Run: func(cmd *cobra.Command, args []string) {
		cli.Println(printablePerennialBranches(prodRepo.GetPerennialBranches()))
	},
	Args: cobra.NoArgs,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return git.ValidateIsRepository()
	},
}

var updatePrennialBranchesCommand = &cobra.Command{
	Use:   "update",
	Short: "Prompts to update your perennial branches",
	Long:  `Prompts to update your perennial branches`,
	Run: func(cmd *cobra.Command, args []string) {
		err := prompt.ConfigurePerennialBranches(prodRepo)
		if err != nil {
			cli.Exit(err)
		}
	},
	Args: cobra.NoArgs,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return git.ValidateIsRepository()
	},
}

// printablePerennialBranches returns a user printable list of perennial branches.
func printablePerennialBranches(perennialBranches []string) string {
	if len(perennialBranches) == 0 {
		return "[none]"
	}
	return strings.Join(perennialBranches, "\n")
}

func init() {
	perennialBranchesCommand.AddCommand(updatePrennialBranchesCommand)
	RootCmd.AddCommand(perennialBranchesCommand)
}
