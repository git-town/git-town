package cmd

import (
	"strings"

	"github.com/git-town/git-town/v7/src/cli"
	"github.com/git-town/git-town/v7/src/validate"
	"github.com/spf13/cobra"
)

const perennialDesc = "Displays your perennial branches"

const perennialHelp = `
Perennial branches are long-lived branches.
They cannot be shipped.`

const updatePerennialSummary = "Prompts to update your perennial branches"

func perennialBranchesCmd() *cobra.Command {
	addDisplayDebugFlag, readDisplayDebugFlag := debugFlag()
	displayCmd := cobra.Command{
		Use:   "perennial-branches",
		Args:  cobra.NoArgs,
		Short: perennialDesc,
		Long:  long(perennialDesc, perennialHelp),
		RunE: func(cmd *cobra.Command, args []string) error {
			return displayPerennialBranches(readDisplayDebugFlag(cmd))
		},
	}
	addDisplayDebugFlag(&displayCmd)

	addUpdateDebugFlag, readUpdateDebugFlag := debugFlag()
	updateCmd := cobra.Command{
		Use:   "update",
		Args:  cobra.NoArgs,
		Short: updatePerennialSummary,
		Long:  long(updatePerennialSummary),
		RunE: func(cmd *cobra.Command, args []string) error {
			return updatePerennialBranches(readUpdateDebugFlag(cmd))
		},
	}
	addUpdateDebugFlag(&updateCmd)
	displayCmd.AddCommand(&updateCmd)
	return &displayCmd
}

func displayPerennialBranches(debug bool) error {
	repo, exit, err := LoadPublicRepo(RepoArgs{
		omitBranchNames:       true,
		debug:                 debug,
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

func updatePerennialBranches(debug bool) error {
	repo, exit, err := LoadPublicRepo(RepoArgs{
		omitBranchNames:       true,
		debug:                 debug,
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
