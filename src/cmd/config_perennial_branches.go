package cmd

import (
	"github.com/git-town/git-town/v9/src/cli/dialog"
	"github.com/git-town/git-town/v9/src/cli/flags"
	"github.com/git-town/git-town/v9/src/cli/format"
	"github.com/git-town/git-town/v9/src/cli/io"
	"github.com/git-town/git-town/v9/src/execute"
	"github.com/spf13/cobra"
)

const perennialDesc = "Displays your perennial branches"

const perennialHelp = `
Perennial branches are long-lived branches.
They cannot be shipped.`

const updatePerennialSummary = "Prompts to update your perennial branches"

func perennialBranchesCmd() *cobra.Command {
	addDisplayDebugFlag, readDisplayDebugFlag := flags.Debug()
	displayCmd := cobra.Command{
		Use:   "perennial-branches",
		Args:  cobra.NoArgs,
		Short: perennialDesc,
		Long:  long(perennialDesc, perennialHelp),
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeConfigPerennialBranches(readDisplayDebugFlag(cmd))
		},
	}
	addDisplayDebugFlag(&displayCmd)

	addUpdateDebugFlag, readUpdateDebugFlag := flags.Debug()
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

func executeConfigPerennialBranches(debug bool) error {
	repo, err := execute.OpenRepo(execute.OpenRepoArgs{
		Debug:            debug,
		DryRun:           false,
		OmitBranchNames:  true,
		ValidateIsOnline: false,
		ValidateGitRepo:  true,
	})
	if err != nil {
		return err
	}
	io.Println(format.StringSetting(repo.Runner.Config.PerennialBranches().Join("\n")))
	return nil
}

func updatePerennialBranches(debug bool) error {
	repo, err := execute.OpenRepo(execute.OpenRepoArgs{
		Debug:            debug,
		DryRun:           false,
		OmitBranchNames:  true,
		ValidateIsOnline: false,
		ValidateGitRepo:  true,
	})
	if err != nil {
		return err
	}
	lineage := repo.Runner.Config.Lineage()
	pushHook, err := repo.Runner.Config.PushHook()
	if err != nil {
		return err
	}
	branches, _, _, exit, err := execute.LoadBranches(execute.LoadBranchesArgs{
		Repo:                  repo,
		Fetch:                 false,
		Debug:                 debug,
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
