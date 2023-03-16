package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v7/src/cli"
	"github.com/git-town/git-town/v7/src/config"
	"github.com/git-town/git-town/v7/src/git"
	"github.com/spf13/cobra"
)

func offlineCmd(repo *git.PublicRepo) *cobra.Command {
	return &cobra.Command{
		Use:   "offline [(yes | no)]",
		Args:  cobra.MaximumNArgs(1),
		Short: "Displays or sets offline mode",
		Long: `Displays or sets offline mode

Git Town avoids network operations in offline mode.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				return setOfflineStatus(args[0], repo)
			}
			return displayOfflineStatus(repo)
		},
	}
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
