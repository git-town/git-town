package cmd

import (
	"github.com/spf13/cobra"
)

const resetConfigSummary = "Resets your Git Town configuration"

func resetConfigCommand() *cobra.Command {
	cmd := cobra.Command{
		Use:   "reset",
		Args:  cobra.NoArgs,
		Short: resetConfigSummary,
		Long:  long(resetConfigSummary),
		RunE:  runConfigReset,
	}
	addDebugFlag(&cmd)
	return &cmd
}

func runConfigReset(cmd *cobra.Command, args []string) error {
	repo, exit, err := LoadPublicRepo(RepoArgs{
		omitBranchNames:       true,
		debug:                 readDebugFlag(cmd),
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
