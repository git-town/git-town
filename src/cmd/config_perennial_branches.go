package cmd

import (
	"strings"

	"github.com/git-town/git-town/v7/src/cli"
	"github.com/git-town/git-town/v7/src/validate"
	"github.com/spf13/cobra"
)

const configPerennialSummary = "Displays your perennial branches"

const configPerennialDesc = `
Perennial branches are long-lived branches.
They cannot be shipped.`

func perennialBranchesCmd() *cobra.Command {
	debug := false
	displayCmd := cobra.Command{
		Use:   "perennial-branches",
		Args:  cobra.NoArgs,
		Short: configPerennialSummary,
		Long:  long(configPerennialSummary, configPerennialDesc),
		RunE:  displayPerennialBranches,
	}
	debugFlagOld(&displayCmd, &debug)
	updateCmd := cobra.Command{
		Use:   "update",
		Short: "Prompts to update your perennial branches",
		Long:  `Prompts to update your perennial branches`,
		RunE:  updatePerennialBranches,
		Args:  cobra.NoArgs,
	}
	debugFlagOld(&updateCmd, &debug)
	displayCmd.AddCommand(&updateCmd)
	return &displayCmd
}

func displayPerennialBranches(cmd *cobra.Command, args []string) error {
	repo, exit, err := LoadPublicRepo(RepoArgs{
		omitBranchNames:       true,
		debug:                 readDebugFlag(cmd),
		dryRun:                false,
		handleUnfinishedState: false,
		validateGitversion:    true,
		validateIsRepository:  true,
	})
	if err != nil || exit {
		return err
	}
	cli.Println(cli.StringSetting(strings.Join(repo.Config.PerennialBranches(), "\n")))
	return nil
}

func updatePerennialBranches(cmd *cobra.Command, args []string) error {
	repo, exit, err := LoadPublicRepo(RepoArgs{
		omitBranchNames:       true,
		debug:                 readDebugFlag(cmd),
		dryRun:                false,
		handleUnfinishedState: false,
		validateGitversion:    true,
		validateIsRepository:  true,
	})
	if err != nil || exit {
		return err
	}
	mainBranch := repo.Config.MainBranch()
	return validate.EnterPerennialBranches(&repo, mainBranch)
}
