package perennialbranches

import (
	"slices"

	"github.com/git-town/git-town/v11/src/cli/dialog/dialogscreens"
	"github.com/git-town/git-town/v11/src/cli/flags"
	"github.com/git-town/git-town/v11/src/cmd/cmdhelpers"
	"github.com/git-town/git-town/v11/src/execute"
	"github.com/spf13/cobra"
)

const updateSummary = "Updates all perennial branches through a visual dialog"

func changeCmd() *cobra.Command {
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	cmd := cobra.Command{
		Use:   "change",
		Args:  cobra.NoArgs,
		Short: updateSummary,
		Long:  cmdhelpers.Long(updateSummary),
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeUpdate(readVerboseFlag(cmd))
		},
	}
	addVerboseFlag(&cmd)
	return &cmd
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
	branchesSnapshot, _, dialogTestInputs, exit, err := execute.LoadRepoSnapshot(execute.LoadBranchesArgs{
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
	newPerennialBranches, aborted, err := dialogscreens.EnterPerennialBranches(branchesSnapshot.Branches.Names(), repo.Runner.PerennialBranches, repo.Runner.MainBranch, dialogTestInputs.Next())
	if err != nil || aborted {
		return err
	}
	if slices.Compare(repo.Runner.PerennialBranches, newPerennialBranches) != 0 {
		return repo.Runner.SetPerennialBranches(newPerennialBranches)
	}
	return nil
}
