package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v7/src/cli"
	"github.com/git-town/git-town/v7/src/config"
	"github.com/git-town/git-town/v7/src/git"
	"github.com/spf13/cobra"
)

const configOfflineSummary = "Displays or sets offline mode"

const configOfflineDesc = `
Git Town avoids network operations in offline mode.`

func offlineCmd() *cobra.Command {
	debug := false
	cmd := cobra.Command{
		Use:   "offline [(yes | no)]",
		Args:  cobra.MaximumNArgs(1),
		Short: configOfflineSummary,
		Long:  long(configOfflineSummary, configOfflineDesc),
		RunE:  configureOffline,
	}
	debugFlagOld(&cmd, &debug)
	return &cmd
}

func configureOffline(cmd *cobra.Command, args []string) error {
	repo, exit, err := LoadPublicRepo(RepoArgs{
		omitBranchNames:       true,
		debug:                 readDebugFlag(cmd),
		dryRun:                false,
		handleUnfinishedState: false,
		validateGitversion:    true,
	})
	if err != nil || exit {
		return err
	}
	if len(args) > 0 {
		return setOfflineStatus(args[0], &repo)
	}
	return displayOfflineStatus(&repo)
}

func displayOfflineStatus(repo *git.PublicRepo) error {
	isOffline, err := repo.Config.IsOffline()
	if err != nil {
		return err
	}
	cli.Println(cli.FormatBool(isOffline))
	return nil
}

func setOfflineStatus(text string, repo *git.PublicRepo) error {
	value, err := config.ParseBool(text)
	if err != nil {
		return fmt.Errorf(`invalid argument: %q. Please provide either "yes" or "no".\n`, text)
	}
	return repo.Config.SetOffline(value)
}
