package config

import (
	"github.com/git-town/git-town/v11/src/cli/dialog"
	"github.com/git-town/git-town/v11/src/cli/flags"
	"github.com/git-town/git-town/v11/src/cmd/cmdhelpers"
	"github.com/git-town/git-town/v11/src/execute"
	"github.com/spf13/cobra"
)

// TODO: extract strings like these into en.go.
const setupConfigDesc = "Prompts to setup your Git Town configuration"

func setupConfigCommand() *cobra.Command {
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	cmd := cobra.Command{
		Use:   "setup",
		Args:  cobra.NoArgs,
		Short: setupConfigDesc,
		Long:  cmdhelpers.Long(setupConfigDesc),
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeConfigSetup(readVerboseFlag(cmd))
		},
	}
	addVerboseFlag(&cmd)
	return &cmd
}

func executeConfigSetup(verbose bool) error {
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
	branches, _, _, exit, err := execute.LoadBranches(execute.LoadBranchesArgs{
		FullConfig:            &repo.Runner.FullConfig,
		Repo:                  repo,
		Verbose:               verbose,
		Fetch:                 false,
		HandleUnfinishedState: false,
		ValidateIsConfigured:  false,
		ValidateNoOpenChanges: false,
	})
	if err != nil || exit {
		return err
	}
	newMainBranch, err := dialog.EnterMainBranch(branches.All.LocalBranches().Names(), branches.Types.MainBranch, &repo.Runner.Backend)
	if err != nil {
		return err
	}
	branches.Types.MainBranch = newMainBranch
	_, err = dialog.EnterPerennialBranches(&repo.Runner.Backend, branches)
	return err
}
