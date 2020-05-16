package cmd

import (
	"github.com/git-town/git-town/src/cfmt"
	"github.com/git-town/git-town/src/git"
	"github.com/git-town/git-town/src/util"
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
			setOfflineFlag(util.StringToBool(args[0]))
		}
	},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 1 {
			return validateBooleanArgument(args[0])
		}
		return cobra.MaximumNArgs(1)(cmd, args)
	},
}

func printOfflineFlag() {
	cfmt.Println(git.GetPrintableOfflineFlag())
}

func setOfflineFlag(value bool) {
	git.Config().SetOffline(value)
}

func init() {
	RootCmd.AddCommand(offlineCommand)
}
