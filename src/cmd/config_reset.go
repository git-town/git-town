package cmd

import (
	"github.com/spf13/cobra"
)

const resetConfigDesc = "Resets your Git Town configuration"

func resetConfigCommand() *cobra.Command {
	addDebugFlag, readDebugFlag := debugFlag()
	cmd := cobra.Command{
		Use:   "reset",
		Args:  cobra.NoArgs,
		Short: resetConfigDesc,
		Long:  long(resetConfigDesc),
		RunE: func(cmd *cobra.Command, args []string) error {
			return resetStatus(readDebugFlag(cmd))
		},
	}
	addDebugFlag(&cmd)
	return &cmd
}

func resetStatus(debug bool) error {
	repo, exit, err := LoadProdRepo(RepoArgs{
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
	return repo.Config.RemoveLocalGitConfiguration()
}
