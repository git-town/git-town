package cmd

import (
	"github.com/git-town/git-town/v11/src/cli/dialog"
	"github.com/git-town/git-town/v11/src/cli/flags"
	"github.com/git-town/git-town/v11/src/cli/format"
	"github.com/git-town/git-town/v11/src/cli/io"
	"github.com/git-town/git-town/v11/src/execute"
	"github.com/spf13/cobra"
)

const perennialDesc = "Displays your perennial branches"

const perennialHelp = `
Perennial branches are long-lived branches.
They cannot be shipped.`

const updatePerennialSummary = "Prompts to update your perennial branches"

func perennialBranchesCmd() *cobra.Command {
	addDisplayVerboseFlag, readDisplayVerboseFlag := flags.Verbose()
	displayCmd := cobra.Command{
		Use:   "perennial-branches",
		Args:  cobra.NoArgs,
		Short: perennialDesc,
		Long:  long(perennialDesc, perennialHelp),
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeConfigPerennialBranches(readDisplayVerboseFlag(cmd))
		},
	}
	addDisplayVerboseFlag(&displayCmd)

	addUpdateVerboseFlag, readUpdateVerboseFlag := flags.Verbose()
	updateCmd := cobra.Command{
		Use:   "update",
		Args:  cobra.NoArgs,
		Short: updatePerennialSummary,
		Long:  long(updatePerennialSummary),
		RunE: func(cmd *cobra.Command, args []string) error {
			return updatePerennialBranches(readUpdateVerboseFlag(cmd))
		},
	}
	addUpdateVerboseFlag(&updateCmd)
	displayCmd.AddCommand(&updateCmd)
	return &displayCmd
}

func executeConfigPerennialBranches(verbose bool) error {
	repo, err := execute.OpenRepo(execute.OpenRepoArgs{
		Verbose:          verbose,
		DryRun:           false,
		OmitBranchNames:  true,
		PrintCommands:    true,
		ValidateIsOnline: false,
		ValidateGitRepo:  true,
	})
	if err != nil {
		return err
	}
	io.Println(format.StringSetting(repo.Runner.GitTown.PerennialBranches().Join("\n")))
	return nil
}

func updatePerennialBranches(verbose bool) error {
	repo, err := execute.OpenRepo(execute.OpenRepoArgs{
		Verbose:          verbose,
		DryRun:           false,
		OmitBranchNames:  true,
		PrintCommands:    true,
		ValidateIsOnline: false,
		ValidateGitRepo:  true,
	})
	if err != nil {
		return err
	}
	lineage := repo.Runner.GitTown.Lineage(repo.Runner.Backend.GitTown.RemoveLocalConfigValue)
	pushHook, err := repo.Runner.GitTown.PushHook()
	if err != nil {
		return err
	}
	branches, _, _, exit, err := execute.LoadBranches(execute.LoadBranchesArgs{
		Repo:                  repo,
		Fetch:                 false,
		Verbose:               verbose,
		HandleUnfinishedState: false,
		Lineage:               lineage,
		PushHook:              pushHook,
		ValidateIsConfigured:  false,
		ValidateNoOpenChanges: false,
	})
	if err != nil || exit {
		return err
	}
	_, err = dialog.EnterPerennialBranches(&repo.Runner.Backend, branches)
	return err
}
