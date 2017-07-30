package cmd

import (
	"fmt"

	"github.com/Originate/git-town/src/git"
	"github.com/Originate/git-town/src/util"
	"github.com/spf13/cobra"
)

var offlineCommand = &cobra.Command{
	Use:   "offline [(true | false)]",
	Short: "Displays or sets offline mode",
	Run: func(cmd *cobra.Command, args []string) {
		git.EnsureIsRepository()
		if len(args) == 0 {
			printOfflineFlag()
		} else {
			setOfflineFlag(util.StringToBool(args[0]))
		}
	},
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 1 {
			err := validateBooleanArgument(args[0])
			if err != nil {
				return err
			}
		}
		return validateMaxArgs(args, 1)
	},
}

func printOfflineFlag() {
	fmt.Println(git.GetPrintableOfflineFlag())
}

func setOfflineFlag(value bool) {
	git.UpdateOffline(value)
}

func init() {
	RootCmd.AddCommand(offlineCommand)
}
