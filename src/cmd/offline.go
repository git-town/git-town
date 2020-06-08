package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/git-town/git-town/src/cli"
	"github.com/git-town/git-town/src/git"
	"github.com/spf13/cobra"
)

var offlineCommand = &cobra.Command{
	Use:   "offline [(true | false)]",
	Short: "Displays or sets offline mode",
	Long: `Displays or sets offline mode

Git Town avoids network operations in offline mode.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			printOfflineFlag()
		} else {
			value, err := strconv.ParseBool(args[0])
			if err != nil {
				fmt.Printf(`Error: invalid argument: %q. Please provide either "true" or "false".\n`, args[0])
				os.Exit(1)
			}
			setOfflineFlag(value)
		}
	},
	Args: cobra.MaximumNArgs(1),
}

func printOfflineFlag() {
	cli.Println(prodRepo.Configuration.PrintableOfflineFlag())
}

func setOfflineFlag(value bool) {
	git.Config().SetOffline(value)
}

func init() {
	RootCmd.AddCommand(offlineCommand)
}
