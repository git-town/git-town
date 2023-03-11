package cmd

import (
	"strings"

	"github.com/git-town/git-town/v7/src/cli"
	"github.com/git-town/git-town/v7/src/dialog"
	"github.com/git-town/git-town/v7/src/git"
	. "github.com/git-town/git-town/v7/src/validate"
	"github.com/spf13/cobra"
)

func perennialBranchesCmd(repo *git.ProdRepo) *cobra.Command {
	perennialBranchesCmd := cobra.Command{
		Use:     "perennial-branches",
		Args:    cobra.NoArgs,
		PreRunE: Ensure(repo, IsRepository),
		Short:   "Displays your perennial branches",
		Long: `Displays your perennial branches

Perennial branches are long-lived branches.
They cannot be shipped.`,
		Run: func(cmd *cobra.Command, args []string) {
			cli.Println(cli.StringSetting(strings.Join(repo.Config.PerennialBranches(), "\n")))
		},
	}
	perennialBranchesCmd.AddCommand(updatePerennialBranchesCmd(repo))
	return &perennialBranchesCmd
}

func updatePerennialBranchesCmd(repo *git.ProdRepo) *cobra.Command {
	return &cobra.Command{
		Use:   "update",
		Short: "Prompts to update your perennial branches",
		Long:  `Prompts to update your perennial branches`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return configurePerennialBranches(repo)
		},
		Args:    cobra.NoArgs,
		PreRunE: Ensure(repo, IsRepository),
	}
}

func configurePerennialBranches(repo *git.ProdRepo) error {
	localBranchesWithoutMain, err := repo.Silent.LocalBranchesWithoutMain()
	if err != nil {
		return err
	}
	perennialBranches := repo.Config.PerennialBranches()
	newPerennialBranches, err := dialog.AskPerennialBranches(localBranchesWithoutMain, perennialBranches)
	if err != nil {
		return err
	}
	return repo.Config.SetPerennialBranches(newPerennialBranches)
}
