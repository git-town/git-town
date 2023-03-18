package cmd

import (
	"strings"

	"github.com/git-town/git-town/v7/src/cli"
	"github.com/git-town/git-town/v7/src/validate"
	"github.com/spf13/cobra"
)

func perennialBranchesCmd() *cobra.Command {
	debug := false
	displayCmd := cobra.Command{
		Use:   "perennial-branches",
		Args:  cobra.NoArgs,
		Short: "Displays your perennial branches",
		Long: `Displays your perennial branches

Perennial branches are long-lived branches.
They cannot be shipped.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return displayPerennialBranches(debug)
		},
	}
	debugFlag(&displayCmd, &debug)
	updateCmd := cobra.Command{
		Use:   "update",
		Short: "Prompts to update your perennial branches",
		Long:  `Prompts to update your perennial branches`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return updatePerennialBranches(debug)
		},
		Args: cobra.NoArgs,
	}
	debugFlag(&updateCmd, &debug)
	displayCmd.AddCommand(&updateCmd)
	return &displayCmd
}

func displayPerennialBranches(debug bool) error {
	repo, err := Repo(RepoArgs{
		omitBranchNames:      true,
		debug:                debug,
		dryRun:               false,
		validateGitversion:   true,
		validateIsRepository: true,
	})
	if err != nil {
		return err
	}
	cli.Println(cli.StringSetting(strings.Join(repo.Config.PerennialBranches(), "\n")))
	return nil
}

func updatePerennialBranches(debug bool) error {
	repo, err := Repo(RepoArgs{
		omitBranchNames:      true,
		debug:                debug,
		dryRun:               false,
		validateGitversion:   true,
		validateIsRepository: true,
	})
	if err != nil {
		return err
	}
	mainBranch := repo.Config.MainBranch()
	return validate.EnterPerennialBranches(&repo, mainBranch)
}
