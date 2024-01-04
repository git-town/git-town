package config

import (
	"fmt"

	"github.com/git-town/git-town/v11/src/cli/dialog"
	"github.com/git-town/git-town/v11/src/cli/flags"
	"github.com/git-town/git-town/v11/src/cli/format"
	"github.com/git-town/git-town/v11/src/cli/io"
	"github.com/git-town/git-town/v11/src/cmd/cmdhelpers"
	"github.com/git-town/git-town/v11/src/execute"
	"github.com/git-town/git-town/v11/src/git/gitdomain"
	"github.com/spf13/cobra"
)

const perennialDesc = "Displays your perennial branches"

const perennialHelp = `
Perennial branches are long-lived branches.
They cannot be shipped.`

const updatePerennialSummary = "Prompts to update your perennial branches"

const addPerennialSummary = "Registers the given branch as a perennial branch"

const removePerennialSummary = "Removes the given branch from the list of perennial branches"

func perennialBranchesCmd() *cobra.Command {
	addDisplayVerboseFlag, readDisplayVerboseFlag := flags.Verbose()
	displayCmd := cobra.Command{
		Use:   "perennial-branches",
		Args:  cobra.NoArgs,
		Short: perennialDesc,
		Long:  cmdhelpers.Long(perennialDesc, perennialHelp),
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
		Long:  cmdhelpers.Long(updatePerennialSummary),
		RunE: func(cmd *cobra.Command, args []string) error {
			return updatePerennialBranches(readUpdateVerboseFlag(cmd))
		},
	}
	addUpdateVerboseFlag(&updateCmd)
	displayCmd.AddCommand(&updateCmd)

	addAddVerboseFlag, readAddVerboseFlag := flags.Verbose()
	addCmd := cobra.Command{
		Use:   "add",
		Args:  cobra.ExactArgs(1),
		Short: addPerennialSummary,
		Long:  cmdhelpers.Long(addPerennialSummary),
		RunE: func(cmd *cobra.Command, args []string) error {
			return addPerennialBranch(readAddVerboseFlag(cmd))
		},
	}
	addAddVerboseFlag(&addCmd)
	displayCmd.AddCommand(&addCmd)

	addRemoveVerboseFlag, readRemoveVerboseFlag := flags.Verbose()
	removeCmd := cobra.Command{
		Use:   "remove",
		Args:  cobra.ExactArgs(1),
		Short: removePerennialSummary,
		Long:  cmdhelpers.Long(removePerennialSummary),
		RunE: func(cmd *cobra.Command, args []string) error {
			return removePerennialBranch(readRemoveVerboseFlag(cmd))
		},
	}
	addRemoveVerboseFlag(&removeCmd)
	displayCmd.AddCommand(&removeCmd)

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
	io.Println(format.StringSetting(repo.Runner.Config.PerennialBranches.Join("\n")))
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
	// check if branch exists
	branchName := gitdomain.NewLocalBranchName(branchStr)
	if !repo.Runner.Backend.HasLocalBranch(branchName) {
		return fmt.Errorf("branch %q does not exist")
	}
	newPerennialBranches := append(repo.Runner.PerennialBranches, branchName)
	return repo.Runner.Config.SetPerennialBranches(newPerennialBranches)
}
