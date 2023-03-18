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
	repo := Repo(debug, false)
	err := ensure(&repo, hasGitVersion, isRepository)
	if err != nil {
		return err
	}
	return repo.Config.RemoveLocalGitConfiguration()
}
