package cmd

import (
	"github.com/spf13/cobra"
)

func resetConfigCommand() *cobra.Command {
	debug := false
	cmd := cobra.Command{
		Use:   "reset",
		Args:  cobra.NoArgs,
		Short: "Resets your Git Town configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runConfigReset(debug)
		},
	}
	debugFlag(&cmd, &debug)
	return &cmd
}

func runConfigReset(debug bool) error {
	repo, err := LoadRepo(RepoArgs{
		omitBranchNames:      true,
		debug:                debug,
		dryRun:               false,
		validateGitversion:   true,
		validateIsRepository: true,
	})
	if err != nil {
		return err
	}
	return repo.Config.RemoveLocalGitConfiguration()
}
