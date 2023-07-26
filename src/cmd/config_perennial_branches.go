package cmd

import (
	"strings"

	"github.com/git-town/git-town/v9/src/cli"
	"github.com/git-town/git-town/v9/src/execute"
	"github.com/git-town/git-town/v9/src/flags"
	"github.com/git-town/git-town/v9/src/validate"
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
			return displayPerennialBranches(readDisplayDebugFlag(cmd))
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

func displayPerennialBranches(debug bool) error {
	run, err := execute.LoadProdRunner(execute.LoadArgs{
		Debug:           debug,
		DryRun:          false,
		OmitBranchNames: true,
	})
	if err != nil {
		return err
	}
	_, _, _, exit, err := execute.LoadGitRepo(&run, execute.LoadGitArgs{
		Fetch:                 false,
		HandleUnfinishedState: false,
		ValidateIsConfigured:  false,
		ValidateIsOnline:      false,
		ValidateNoOpenChanges: false,
	})
	if err != nil || exit {
		return err
	}
	cli.Println(cli.StringSetting(strings.Join(run.Config.PerennialBranches(), "\n")))
	return nil
}

func updatePerennialBranches(debug bool) error {
	run, err := execute.LoadProdRunner(execute.LoadArgs{
		Debug:           debug,
		DryRun:          false,
		OmitBranchNames: true,
	})
	if err != nil {
		return err
	}
	allBranches, _, _, exit, err := execute.LoadGitRepo(&run, execute.LoadGitArgs{
		Fetch:                 false,
		HandleUnfinishedState: false,
		ValidateIsConfigured:  false,
		ValidateIsOnline:      false,
		ValidateNoOpenChanges: false,
	})
	if err != nil || exit {
		return err
	}
	mainBranch := run.Config.MainBranch()
	return validate.EnterPerennialBranches(&run.Backend, allBranches, mainBranch)
}
