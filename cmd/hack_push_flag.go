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
			printHackPushFlag()
		} else {
			setHackPushFlag(stringToBool(args[0]))
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

func printHackPushFlag() {
	fmt.Println(strconv.FormatBool(git.ShouldHackPush()))
}

func setHackPushFlag(value bool) {
	git.UpdateShouldHackPush(value)
}

func init() {
	RootCmd.AddCommand(hackPushFlagCommand)
}
