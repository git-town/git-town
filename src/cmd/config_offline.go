package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v7/src/cli"
	"github.com/git-town/git-town/v7/src/config"
	"github.com/git-town/git-town/v7/src/git"
	"github.com/spf13/cobra"
)

const offlineDesc = "Displays or sets offline mode"

const offlineHelp = `
Git Town avoids network operations in offline mode.`

func offlineCmd(repo *git.ProdRepo) *cobra.Command {
	return &cobra.Command{
		Use:   "offline [(yes | no)]",
		Args:  cobra.MaximumNArgs(1),
		Short: offlineDesc,
		Long:  long(offlineDesc, offlineHelp),
		RunE: func(cmd *cobra.Command, args []string) error {
			return offline(args, repo)
		},
	}
}

func offline(args []string, repo *git.ProdRepo) error {
	if len(args) > 0 {
		return setOfflineStatus(args[0], repo)
	}
	return displayOfflineStatus(repo)
}

func displayOfflineStatus(repo *git.ProdRepo) error {
	isOffline, err := repo.Config.IsOffline()
	if err != nil {
		return err
	}
	cli.Println(cli.FormatBool(isOffline))
	return nil
}

func setOfflineStatus(text string, repo *git.ProdRepo) error {
	value, err := config.ParseBool(text)
	if err != nil {
		return fmt.Errorf(`invalid argument: %q. Please provide either "yes" or "no".\n`, text)
	}
	return repo.Config.SetOffline(value)
}
