package cmd

import (
	"fmt"
	"strconv"

	"github.com/Originate/git-town/lib/git"
	"github.com/spf13/cobra"
)

var hackPushFlagCommand = &cobra.Command{
	Use:   "hack-push-flag [(true | false)]",
	Short: "Displays or sets your hack push flag",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println(strconv.FormatBool(git.ShouldHackPush()))
			return
		}

		git.UpdateShouldHackPush(args[0] == "true")
	},
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 1 && args[0] != "true" && args[0] != "false" {
			return fmt.Errorf("Invalid value: '%s'", args[0])
		}
		return validateMaxArgs(args, 1)
	},
}

func init() {
	RootCmd.AddCommand(hackPushFlagCommand)
}
