package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v7/src/cli"
	"github.com/git-town/git-town/v7/src/git"
	"github.com/spf13/cobra"
)

func offlineCmd(repo *git.ProdRepo) *cobra.Command {
	return &cobra.Command{
		Use:   "offline [(yes | no)]",
		Short: "Displays or sets offline mode",
		Long: `Displays or sets offline mode

Git Town avoids network operations in offline mode.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				isOffline, err := repo.Config.IsOffline()
				if err != nil {
					return err
				}
				cli.Println(cli.FormatBool(isOffline))
			} else {
				value, err := cli.ParseBool(args[0])
				if err != nil {
					return fmt.Errorf(`invalid argument: %q. Please provide either "yes" or "no".\n`, args[0])
				}
				err = repo.Config.SetOffline(value)
				if err != nil {
					return err
				}
			}
			return nil
		},
		Args: cobra.MaximumNArgs(1),
	}
}
