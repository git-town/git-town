package cmd

import (
	"github.com/git-town/git-town/v7/src/validate"
	"github.com/spf13/cobra"
)

const setupConfigSummary = "Prompts to setup your Git Town configuration"

func setupConfigCommand() *cobra.Command {
	cmd := cobra.Command{
		Use:   "setup",
		Args:  cobra.NoArgs,
		Short: setupConfigSummary,
		Long:  long(setupConfigSummary),
		RunE:  runConfigSetup,
	}
	addDebugFlag(&cmd)
	return &cmd
}

func runConfigSetup(cmd *cobra.Command, args []string) error {
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
	mainBranch, err := validate.EnterMainBranch(&repo)
	if err != nil {
		return err
	}
	return validate.EnterPerennialBranches(&repo, mainBranch)
}
