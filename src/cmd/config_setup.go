package cmd

import (
	"github.com/git-town/git-town/v7/src/validate"
	"github.com/spf13/cobra"
)

const setupConfigDesc = "Prompts to setup your Git Town configuration"

func setupConfigCommand() *cobra.Command {
	addDebugFlag, readDebugFlag := debugFlag()
	cmd := cobra.Command{
		Use:   "setup",
		Args:  cobra.NoArgs,
		Short: setupConfigDesc,
		Long:  long(setupConfigDesc),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runConfigSetup(readDebugFlag(cmd))
		},
	}
	addDebugFlag(&cmd)
	return &cmd
}

func runConfigSetup(debug bool) error {
	repo, exit, err := LoadPublicRepo(RepoArgs{
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
	mainBranch, err := validate.EnterMainBranch(&repo)
	if err != nil {
		return err
	}
	return validate.EnterPerennialBranches(&repo, mainBranch)
}
