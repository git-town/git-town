package perennialbranches

import (
	"github.com/git-town/git-town/v11/src/cli/dialog"
	"github.com/git-town/git-town/v11/src/cli/flags"
	"github.com/git-town/git-town/v11/src/cmd/cmdhelpers"
	"github.com/git-town/git-town/v11/src/execute"
	"github.com/spf13/cobra"
)

const updatePerennialSummary = "Prompts to update your perennial branches"

func updateCmd() *cobra.Command {
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	updateCmd := cobra.Command{
		Use:   "update",
		Args:  cobra.NoArgs,
		Short: updatePerennialSummary,
		Long:  cmdhelpers.Long(updatePerennialSummary),
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeUpdate(readVerboseFlag(cmd))
		},
	}
	addVerboseFlag(&updateCmd)
	return &updateCmd
}

func executeUpdate(verbose bool) error {
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
	branchesSnapshot, _, exit, err := execute.LoadRepoSnapshot(execute.LoadBranchesArgs{
		FullConfig:            &repo.Runner.FullConfig,
		Repo:                  repo,
		Fetch:                 false,
		Verbose:               verbose,
		HandleUnfinishedState: false,
		ValidateIsConfigured:  false,
		ValidateNoOpenChanges: false,
	})
	if err != nil || exit {
		return err
	}
	err = dialog.EnterPerennialBranches(&repo.Runner.Backend, &repo.Runner.FullConfig, branchesSnapshot.Branches)
	return err
}
