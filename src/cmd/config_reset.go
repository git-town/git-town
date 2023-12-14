package cmd

import (
	"github.com/git-town/git-town/v11/src/cli/flags"
	"github.com/git-town/git-town/v11/src/execute"
	"github.com/spf13/cobra"
)

const resetConfigDesc = "Resets your Git Town configuration"

func resetConfigCommand() *cobra.Command {
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	cmd := cobra.Command{
		Use:   "reset",
		Args:  cobra.NoArgs,
		Short: resetConfigDesc,
		Long:  long(resetConfigDesc),
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeConfigResetStatus(readVerboseFlag(cmd))
		},
	}
	addVerboseFlag(&cmd)
	return &cmd
}

func executeConfigResetStatus(verbose bool) error {
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
	return repo.Runner.GitTown.RemoveLocalGitConfiguration()
}
