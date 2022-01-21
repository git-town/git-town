package cmd

import (
	"fmt"
	"strconv"

	"github.com/git-town/git-town/v7/src/cli"
	"github.com/spf13/cobra"
)

var offlineCommand = &cobra.Command{
	Use:   "offline [(true | false)]",
	Short: "Displays or sets offline mode",
	Long: `Displays or sets offline mode

Git Town avoids network operations in offline mode.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cli.Println(cli.PrintableOfflineFlag(prodRepo.Config.IsOffline()))
		} else {
			value, err := strconv.ParseBool(args[0])
			if err != nil {
				cli.Exit(fmt.Errorf(`invalid argument: %q. Please provide either "true" or "false".\n`, args[0]))
			}
			err = prodRepo.Config.SetOffline(value)
			if err != nil {
				cli.Exit(err)
			}
		}
	},
	Args: cobra.MaximumNArgs(1),
}

func init() {
	RootCmd.AddCommand(offlineCommand)
}
