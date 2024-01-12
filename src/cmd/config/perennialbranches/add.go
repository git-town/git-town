package perennialbranches

import (
	"fmt"
	"slices"

	"github.com/git-town/git-town/v11/src/cli/flags"
	"github.com/git-town/git-town/v11/src/cmd/cmdhelpers"
	"github.com/git-town/git-town/v11/src/execute"
	"github.com/git-town/git-town/v11/src/git/gitdomain"
	"github.com/spf13/cobra"
)

const addPerennialSummary = "Registers the given branch as a perennial branch"

func addCmd() *cobra.Command {
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	cmd := cobra.Command{
		Use:   "add",
		Args:  cobra.ExactArgs(1),
		Short: addPerennialSummary,
		Long:  cmdhelpers.Long(addPerennialSummary),
		RunE: func(cmd *cobra.Command, args []string) error {
			return addPerennialBranch(args[0], readVerboseFlag(cmd))
		},
	}
	addVerboseFlag(&cmd)
	return &cmd
}

func addPerennialBranch(branchStr string, verbose bool) error {
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
	_, _, _, exit, err := execute.LoadRepoSnapshot(execute.LoadBranchesArgs{
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
	branchName := gitdomain.NewLocalBranchName(branchStr)
	if !repo.Runner.Backend.HasLocalBranch(branchName) {
		return fmt.Errorf("branch %q does not exist", branchName)
	}
	if slices.Contains(repo.Runner.PerennialBranches, branchName) {
		return fmt.Errorf("branch %q is already perennial", branchName)
	}
	newPerennialBranches := append(repo.Runner.PerennialBranches, branchName) //nolint:gocritic
	newPerennialBranches.Sort()
	return repo.Runner.Config.SetPerennialBranches(newPerennialBranches)
}
