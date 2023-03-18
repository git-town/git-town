package cmd

import (
	"github.com/git-town/git-town/v7/src/validate"
	"github.com/spf13/cobra"
)

func setupConfigCommand() *cobra.Command {
	debug := false
	cmd := cobra.Command{
		Use:   "setup",
		Args:  cobra.NoArgs,
		Short: "Prompts to setup your Git Town configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runConfigSetup(debug)
		},
	}
	debugFlag(&cmd, &debug)
	return &cmd
}

func runConfigSetup(debug bool) error {
	repo, err := LoadPublicRepo(RepoArgs{
		omitBranchNames:      true,
		debug:                debug,
		dryRun:               false,
		validateGitversion:   true,
		validateIsRepository: true,
	})
	if err != nil {
		return err
	}
	mainBranch, err := validate.EnterMainBranch(&repo)
	if err != nil {
		return err
	}
	return validate.EnterPerennialBranches(&repo, mainBranch)
}
